package rpc

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"

type WeCrossRPCFactory struct{}

func (w WeCrossRPCFactory) Build(weCrossService service.WeCrossService) (WeCrossRPC, error) {
	err := weCrossService.InitService()
	return &WeCrossRPCRest{weCrossService}, err
}
