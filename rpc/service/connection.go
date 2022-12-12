package service

type Connection struct {
	Server      string
	SslKey      string
	SslCert     string
	CaCert      string
	MaxTotal    int
	MaxPerRoute int
}
