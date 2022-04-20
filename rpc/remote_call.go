package rpc

type RemoteCall struct {
	weCrossService WeCrossService
	httpMethod     string
	uri            string
	responseType   Response
	request        Request
	xx             interface{}
}

func (r *RemoteCall) Send() (Response, error) {
	return nil, nil
}

func (r *RemoteCall) AsyncSend(callback) {
}
