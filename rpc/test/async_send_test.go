package test

import (
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
)

func TestAsyncSend(t *testing.T) {
	weCrossRPCService := &service.WeCrossRPCService{}
	weCrossRPC, err := rpc.WeCrossRPCFactory{}.Build(weCrossRPCService)
	if err != nil {
		t.Fatal(err)
	}
	weCrossRPC.Test().AsyncSend(
		func() *methods.Callback {
			bc := methods.CallbackFactory{}.Build()
			bc.OnSuccess(
				func(response methods.Response) {
					println(response.GetVersion())
				})
			bc.OnFailed(
				func(e errors.Error) {
					println(e.Detail)
				})
			return bc
		}(),
	)
}
