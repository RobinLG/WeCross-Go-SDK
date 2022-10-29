package methods

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"

type Request struct {
	version  string
	data     interface{}
	ext      interface{}            `json:"-"`
	callback common.WeCrossCallback `json:"-"`
}
