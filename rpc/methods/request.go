package methods

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"

type Request struct {
	Version  string
	Data     interface{}
	Ext      interface{}            `json:"-"`
	Callback common.WeCrossCallback `json:"-"`
}
