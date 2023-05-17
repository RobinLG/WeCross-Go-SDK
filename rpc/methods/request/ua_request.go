package request

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common/account"

type UARequest struct {
	username     string
	password     string
	clientType   string
	authCode     string
	chainAccount account.ChainAccount
}

func (u *UARequest) GetUsername() string {
	return u.username
}

func (u *UARequest) SetUsername(username string) {
	u.username = username
}

func (u *UARequest) GetPassword() string {
	return u.password
}

func (u *UARequest) SetPassword(password string) {
	u.password = password
}

func (u *UARequest) GetClientType() string {
	return u.clientType
}

func (u *UARequest) SetClientType(clientType string) {
	u.clientType = clientType
}

func (u *UARequest) GetAuthCode() string {
	return u.authCode
}

func (u *UARequest) SetAuthCode(authCode string) {
	u.authCode = authCode
}

func (u *UARequest) GetChainAccount() account.ChainAccount {
	return u.chainAccount
}

func (u *UARequest) SetChainAccount(chainAccount account.ChainAccount) {
	u.chainAccount = chainAccount
}
