package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
)

type WeCrossRPCRest struct {
	WecrossService service.WeCrossService
}

func (w *WeCrossRPCRest) Test() *RemoteCall {
	return &RemoteCall{
		WeCrossService: w.WecrossService, HttpMethod: "POST", Uri: "/sys/test", Response: &methods.UnimplementedResponse{}, Request: &methods.Request{Version: common.CURRENT_VERSION},
	}
}
