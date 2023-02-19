package jsonrpc

import "encoding/json"

type Handler interface {
	Shutdown() bool
	OnRequest(method string, id RequestID, params json.RawMessage)
	OnNotification(method string, params json.RawMessage)
	OnResult(id RequestID, result json.RawMessage, err *Error)
	OnError(err error)
}
