package server

import (
	"encoding/json"
	"sync"

	"github.com/trwk76/golsp/jsonrpc"
)

type clientRequests struct {
	mtx    sync.Mutex
	ids    map[jsonrpc.RequestID]*clientRequest
	nextID uint64
}

func newClientRequests() clientRequests {
	return clientRequests{
		mtx:    sync.Mutex{},
		ids:    make(map[jsonrpc.RequestID]*clientRequest),
		nextID: 0,
	}
}

func (c *clientRequests) open(wg *sync.WaitGroup) *clientRequest {
	c.mtx.Lock()
	id := c.nextID
	c.nextID += 1

	res := &clientRequest{
		wg: wg,
		id: jsonrpc.NewIntRequestID(id),
	}
	c.ids[res.id] = res

	c.mtx.Unlock()

	if wg != nil {
		wg.Add(1)
	}

	return res
}

func (c *clientRequests) close(id jsonrpc.RequestID, result json.RawMessage, err *jsonrpc.Error) {
	c.mtx.Lock()

	req, ok := c.ids[id]
	if !ok {
		c.mtx.Unlock()
		return
	}

	delete(c.ids, id)
	c.mtx.Unlock()

	req.res = result
	req.err = err

	if req.wg != nil {
		req.wg.Done()
	}
}

type clientRequest struct {
	id  jsonrpc.RequestID
	wg  *sync.WaitGroup
	res json.RawMessage
	err *jsonrpc.Error
}
