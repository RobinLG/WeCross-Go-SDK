package errors

import "encoding/json"

const (
	Success int32 = 0

	InternalError int32 = 100

	// config
	FieldMissing  int32 = 101
	ResourceError int32 = 102
	IllegalSymbol int32 = 103

	// rpc
	RemoteCallError    int32 = 201
	RpcError           int32 = 202
	CallContractError  int32 = 203
	LackAuthentication int32 = 204

	// performance
	ResourceInactive int32 = 301
	InvalidContract  int32 = 302
)

type Error struct {
	Code   int32  `json:"code"`
	Detail string `json:"detail"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}

func (e *Error) ErrCode() int32 {
	return e.Code
}

func New(detail string, code int32) error {
	return &Error{
		Code:   code,
		Detail: detail,
	}
}
