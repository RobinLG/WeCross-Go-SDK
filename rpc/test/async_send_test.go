package test

import (
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosstest"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
)

type s struct {
	wecrosstest.Tester
}

func Test(t *testing.T) {
	wecrosstest.RunSubTests(t, s{})
}

func (s) TestAsyncSend(t *testing.T) {
	weCrossRPCService := &service.WeCrossRPCService{}
	weCrossRPC, err := rpc.WeCrossRPCFactory{}.Build(weCrossRPCService)
	if err.Code != 0 {
		t.Fatal(err)
	}
	weCrossRPC.Test().AsyncSend(
		func() *methods.Callback {
			bc := methods.CallbackFactory{}.Build()
			bc.OnSuccess = func(response methods.Response) {
				println(response.GetVersion())
			}
			bc.OnFailed = func(e *errors.Error) {
				println(e.Detail)
			}
			return bc
		}(),
	)
}
