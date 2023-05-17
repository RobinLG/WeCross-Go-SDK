package methods

//type Response interface {
//	GetVersion() string
//	GetErrorCode() int32
//	GetMessage() string
//	GetData() interface{}
//}

//type UnimplementedResponse struct {
//	version   string
//	errorCode int32
//	message   string
//	data      interface{}
//}
//
//func (UnimplementedResponse) GetVersion() string {
//	return "0.0.0"
//}
//
//func (UnimplementedResponse) GetErrorCode() int32 {
//	return errors.InternalError
//}
//
//func (UnimplementedResponse) GetMessage() string {
//	return ""
//}
//
//func (UnimplementedResponse) GetData() interface{} {
//	return nil
//}

type Response struct {
	version   string
	errorCode int32
	message   string
	data      interface{}
}

func (r *Response) GetVersion() string {
	return r.version
}

func (r *Response) SetVersion(version string) {
	r.version = version
}

func (r *Response) GetErrorCode() int32 {
	return r.errorCode
}

func (r *Response) SetErrorCode(errorCode int32) {
	r.errorCode = errorCode
}

func (r *Response) GetMessage() string {
	return r.message
}

func (r *Response) SetMessage(message string) {
	r.message = message
}

func (r *Response) GetData() interface{} {
	return r.data
}

func (r *Response) SetData(data interface{}) {
	r.data = data
}
