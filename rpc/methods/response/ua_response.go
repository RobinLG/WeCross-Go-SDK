package response

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
)

type UAResponse struct {
	methods.Response
}

func (u *UAResponse) GetUAReceipt() *common.UAReceipt {
	uar, ok := u.GetData().(*common.UAReceipt)
	if ok {
		return uar
	} else {
		return nil
	}
}

func (u *UAResponse) SetUAReceipt(uaLoginReceipt *common.UAReceipt) {
	u.SetData(uaLoginReceipt)
}
