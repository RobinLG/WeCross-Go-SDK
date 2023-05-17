package common

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common/account"

type UAReceipt struct {
	errorCode        int32
	message          string
	credential       string
	universalAccount account.UniversalAccount
}

func (u *UAReceipt) GetErrorCode() int32 {
	return u.errorCode
}

func (u *UAReceipt) SetErrorCode(errorCode int32) {
	u.errorCode = errorCode
}

func (u *UAReceipt) GetMessage() string {
	return u.message
}

func (u *UAReceipt) SetMessage(message string) {
	u.message = message
}

func (u *UAReceipt) GetCredential() string {
	return u.credential
}

func (u *UAReceipt) SetCredential(credential string) {
	u.credential = credential
}

func (u *UAReceipt) GetUniversalAccount() account.UniversalAccount {
	return u.universalAccount
}

func (u *UAReceipt) SetUniversalAccount(universalAccount account.UniversalAccount) {
	u.universalAccount = universalAccount
}
