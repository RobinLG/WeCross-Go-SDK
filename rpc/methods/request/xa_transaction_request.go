package request

type XATransactionRequest struct {
	xaTransactionID string
	paths           []string
}

func (x *XATransactionRequest) GetXaTransactionID() string {
	return x.xaTransactionID
}

func (x *XATransactionRequest) SetXaTransactionID(xaTransactionID string) {
	x.xaTransactionID = xaTransactionID
}

func (x *XATransactionRequest) GetPaths() []string {
	return x.paths
}

func (x *XATransactionRequest) SetPaths(paths []string) {
	x.paths = paths
}
