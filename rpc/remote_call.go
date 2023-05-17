package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
)

type RemoteCall struct {
	WeCrossService service.WeCrossService
	HttpMethod     string
	Uri            string
	Response       methods.Response
	Request        *methods.Request
}

func (r *RemoteCall) Send() (methods.Response, error) {
	return r.WeCrossService.Send(r.HttpMethod, r.Uri, r.Request, r.Response)
}

func (r *RemoteCall) AsyncSend(callback *methods.Callback) {
	r.WeCrossService.AsyncSend(r.HttpMethod, r.Uri, r.Request, r.Response, callback)
}
