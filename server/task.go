package server

import "sync"

type taskQueue struct {
	wg   sync.WaitGroup
	mtx  sync.Mutex
	fifo []task
}

func newTaskQueue() taskQueue {
	return taskQueue{
		wg:   sync.WaitGroup{},
		mtx:  sync.Mutex{},
		fifo: make([]task, 0, 10),
	}
}

func (q *taskQueue) take() task {
	q.wg.Wait()

	q.mtx.Lock()
	empty := (len(q.fifo) == 1)
	res := q.fifo[0]
	q.fifo = q.fifo[1:]
	q.mtx.Unlock()

	if empty {
		q.wg.Add(1)
	}

	return res
}

func (q *taskQueue) append(t task) {
	q.mtx.Lock()

	empty := (len(q.fifo) == 0)
	q.fifo = append(q.fifo, t)

	q.mtx.Unlock()

	if empty {
		q.wg.Done()
	}
}

type task interface {
}
