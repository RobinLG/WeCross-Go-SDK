package service

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"

type WeCrossService interface {
	InitService() error
	Send(httpMethod string, uri string, request *methods.Request, responseType methods.Response) (methods.Response, error)
	AsyncSend(httpMethod string, uri string, request *methods.Request, responseType methods.Response, callback *methods.Callback)
}
