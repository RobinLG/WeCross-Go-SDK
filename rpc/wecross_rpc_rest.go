package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("rpc_rest")

type WeCrossRPCRest struct {
	WecrossService service.WeCrossService
}

func (w *WeCrossRPCRest) Test() *service.RemoteCall {
	return &service.RemoteCall{
		WeCrossService: w.WecrossService, HttpMethod: "POST", Uri: "/sys/test", Response: &methods.UnimplementedResponse{}, Request: &methods.Request{Version: common.CURRENT_VERSION},
	}
}
