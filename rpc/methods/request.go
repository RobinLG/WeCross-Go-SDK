package methods

import "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/common"

type Request struct {
	version  string
	data     interface{}
	ext      interface{}
	callback *common.WeCrossCallback
}

func (r *Request) GetData() interface{} {
	return r.data
}

func (r *Request) SetData(data interface{}) {
	r.data = data
}

func (r *Request) GetExt() interface{} {
	return r.ext
}

func (r *Request) SetExt(ext interface{}) {
	r.ext = ext
}

func (r *Request) GetCallback() *common.WeCrossCallback {
	return r.callback
}

func (r *Request) SetCallback(callback *common.WeCrossCallback) {
	r.callback = callback
}

func (r *Request) GetVersion() string {
	return r.version
}

func (r *Request) SetVersion(version string) {
	r.version = version
}
