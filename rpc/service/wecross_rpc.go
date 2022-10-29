package service

import (
	"net/http"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("wecross_rpc")

type WeCrossRPCService struct {
	server     string
	httpClient http.Client
	urlPrefix  string
}

func (w *WeCrossRPCService) InitService() error {
	return nil
}

func (w *WeCrossRPCService) Send(httpMethod string, uri string, request *methods.Request, responseType methods.Response) (methods.Response, error) {
	return nil, nil
}

func (w *WeCrossRPCService) AsyncSend(httpMethod string, uri string, request *methods.Request, responseType methods.Response, callback *methods.Callback) {
}
