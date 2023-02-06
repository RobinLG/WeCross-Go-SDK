package service

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
)

type WeCrossService interface {
	InitService() *errors.Error
	Send(httpMethod string, uri string, request *methods.Request, responseType methods.Response) (methods.Response, error)
	AsyncSend(httpMethod string, uri string, request *methods.Request, responseType methods.Response, callback *methods.Callback)
}
