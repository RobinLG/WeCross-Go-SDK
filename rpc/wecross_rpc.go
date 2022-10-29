package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
)

type WeCrossRPC interface {
	Test() *service.RemoteCall
}
