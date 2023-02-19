package jsonrpc

import "encoding/json"

type Error struct {
	Code    ErrorCode       `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (e Error) Error() string {
	return e.Message
}

type ErrorCode int32

const (
	ErrorCode_ParseError           ErrorCode = -32700
	ErrorCode_InvalidRequest       ErrorCode = -32600
	ErrorCode_MethodNotFound       ErrorCode = -32601
	ErrorCode_InvalidParams        ErrorCode = -32602
	ErrorCode_InternalError        ErrorCode = -32603
	ErrorCode_ServerNotInitialized ErrorCode = -32002
	ErrorCode_UnknownErrorCode     ErrorCode = -32001
	ErrorCode_RequestFailed        ErrorCode = -32803
	ErrorCode_ServerCancelled      ErrorCode = -32802
	ErrorCode_ContentModified      ErrorCode = -32801
	ErrorCode_RequestCancelled     ErrorCode = -32800
)

var (
	ErrParseError Error = Error{
		Code:    ErrorCode_ParseError,
		Message: "Parse error",
	}

	ErrInvalidRequest Error = Error{
		Code:    ErrorCode_InvalidRequest,
		Message: "Invalid Request",
	}

	ErrMethodNotFound Error = Error{
		Code:    ErrorCode_MethodNotFound,
		Message: "Method not found",
	}

	ErrInvalidParams Error = Error{
		Code:    ErrorCode_InvalidParams,
		Message: "Invalid params",
	}

	ErrInternalError Error = Error{
		Code:    ErrorCode_InternalError,
		Message: "Internal error",
	}

	ErrServerNotInitialized Error = Error{
		Code:    ErrorCode_ServerNotInitialized,
		Message: "Server not initialized",
	}

	ErrUnknownErrorCode Error = Error{
		Code:    ErrorCode_UnknownErrorCode,
		Message: "Unknown error code",
	}

	ErrRequestFailed Error = Error{
		Code:    ErrorCode_RequestFailed,
		Message: "Request failed",
	}

	ErrServerCancelled Error = Error{
		Code:    ErrorCode_ServerCancelled,
		Message: "Server cancelled",
	}

	ErrContentModified Error = Error{
		Code:    ErrorCode_ContentModified,
		Message: "Content modified",
	}

	ErrRequestCancelled Error = Error{
		Code:    ErrorCode_RequestCancelled,
		Message: "Request cancelled",
	}
)
