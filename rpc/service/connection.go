package service

type Connection struct {
	Server    string
	SslKey    string
	SslCert   string
	CaCert    string
	SslSwitch int
	UrlPrefix string
}
