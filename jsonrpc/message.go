package jsonrpc

import "encoding/json"

type Message struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      RequestID       `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

func (m Message) request() bool {
	return (m.Method != "") && (m.ID != "")
}

func (m Message) notification() bool {
	return (m.Method != "") && (m.ID == "")
}

func (m Message) result() bool {
	return (m.Method == "") && (m.ID != "")
}

const Version string = "2.0"
