package response

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"

type RawXAResponse struct {
	status             int32
	chainErrorMessages []common.ChainErrorMessage
}

func (r *RawXAResponse) GetStatus() int32 {
	return r.status
}

func (r *RawXAResponse) SetStatus(status int32) {
	r.status = status
}

func (r *RawXAResponse) getChainErrorMessages() []common.ChainErrorMessage {
	return r.chainErrorMessages
}

func (r *RawXAResponse) setChainErrorMessages(chainErrorMessages []common.ChainErrorMessage) {
	r.chainErrorMessages = chainErrorMessages
}
