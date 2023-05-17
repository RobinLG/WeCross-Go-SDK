package account

type UniversalAccount struct {
	username      string
	password      string
	pubKey        string
	secKey        string
	uaID          string
	chainAccounts []ChainAccount
}
