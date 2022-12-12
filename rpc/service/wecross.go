package service

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"

type WeCrossService interface {
	InitService() error
	Send(request *methods.Request, responseType methods.Response) (methods.Response, error)
	AsyncSend(request *methods.Request, responseType methods.Response, callback *methods.Callback)
}
