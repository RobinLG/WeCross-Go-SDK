package response

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
)

type XAResponse struct {
	methods.Response
}

func (x *XAResponse) GetXARawResponse() *RawXAResponse {
	uar, ok := x.GetData().(*RawXAResponse)
	if ok {
		return uar
	} else {
		return nil
	}
}

func (x *XAResponse) SetXARawResponse(res *RawXAResponse) {
	x.SetData(res)
}
