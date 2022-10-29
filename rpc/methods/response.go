package methods

import "github.com/WeBankBlockchain/WeCross-Go-SDK/errors"

type Response interface {
	GetVersion() string
	GetErrorCode() int32
	GetMessage() string
	GetData() interface{}
}

type UnimplementedResponse struct {
	version   string
	errorCode int32
	message   string
	data      interface{}
}

func (UnimplementedResponse) GetVersion() string {
	return "0.0.0"
}

func (UnimplementedResponse) GetErrorCode() int32 {
	return errors.InternalError
}

func (UnimplementedResponse) GetMessage() string {
	return ""
}

func (UnimplementedResponse) GetData() interface{} {
	return nil
}
