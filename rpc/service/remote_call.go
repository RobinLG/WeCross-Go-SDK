package service

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
)

type RemoteCall struct {
	WeCrossService WeCrossService
	Response       methods.Response
	Request        *methods.Request
}

func (r *RemoteCall) Send() (methods.Response, error) {
	return r.WeCrossService.Send(r.Request, r.Response)
}

func (r *RemoteCall) AsyncSend(callback *methods.Callback) {
	r.WeCrossService.AsyncSend(r.Request, r.Response, callback)
}
