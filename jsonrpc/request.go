package jsonrpc

import (
	"encoding/json"
)

type RequestID string

func NewIntRequestID(value uint64) RequestID {
	raw, _ := json.Marshal(&value)
	return RequestID(string(raw))
}

func NewStringRequestID(value string) RequestID {
	raw, _ := json.Marshal(&value)
	return RequestID(string(raw))
}

func (r *RequestID) UnmarshalJSON(data []byte) error {
	(*r) = RequestID(data)
	return nil
}

func (r *RequestID) MarshalJSON() ([]byte, error) {
	return []byte(*r), nil
}

const NullRequestID RequestID = RequestID("null")
