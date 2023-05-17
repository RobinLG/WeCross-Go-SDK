package rpc

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("rpc")

type WeCrossRPC interface {
	Test() *RemoteCall
}
