package server

import (
	"encoding/json"
	"sync"

	"github.com/trwk76/golsp/jsonrpc"
)

type Session struct {
	c     jsonrpc.Connection
	i     Implementation
	s     SessionState
	creqs clientRequests
	tasks taskQueue
}

func NewSession(c jsonrpc.Connection, impl Implementation) *Session {
	return &Session{
		c:     c,
		i:     impl,
		s:     SessionState_Uninitialized,
		creqs: newClientRequests(),
		tasks: newTaskQueue(),
	}
}

func RequestClient[PA any, RE any](s *Session, method string, params PA) (*RE, *jsonrpc.Error) {
	var res RE

	raw, err := json.Marshal(&params)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	creq := s.creqs.open(&wg)

	s.c.Send(jsonrpc.Message{
		JsonRPC: jsonrpc.Version,
		Method:  method,
		ID:      creq.id,
		Params:  raw,
	})

	wg.Wait()

	if creq.err != nil {
		return nil, creq.err
	}

	if err = json.Unmarshal(creq.res, &res); err != nil {
		panic(err)
	}

	return &res, nil
}

func (s *Session) OnResult(id jsonrpc.RequestID, result json.RawMessage, err *jsonrpc.Error) {
	s.creqs.close(id, result, err)
}

type SessionState string

const (
	SessionState_Uninitialized SessionState = "uninitialized"
	SessionState_Initialized   SessionState = "initialized"
	SessionState_Shutdown      SessionState = "shutdown"
)
