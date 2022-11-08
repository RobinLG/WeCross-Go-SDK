package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/util"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("wecross_rpc")

type WeCrossRPCService struct {
	server     string
	httpClient http.Client
	urlPrefix  string
}

func (w *WeCrossRPCService) InitService() error {
	return nil
}

func (w *WeCrossRPCService) Send(httpMethod string, uri string, request *methods.Request, responseType methods.Response) (methods.Response, error) {
	return nil, nil
}

func (w *WeCrossRPCService) AsyncSend(httpMethod string, uri string, request *methods.Request, responseType methods.Response, callback *methods.Callback) {
	defer util.RecoverError(callback)
	url := ""
	if w.urlPrefix != "" {
		url = fmt.Sprintf("%s%s%s", w.server, w.urlPrefix, uri)
	} else {
		url = fmt.Sprintf("%s%s", w.server, uri)
	}

	checkErr := w.checkRequest(request)
	if checkErr.Code != errors.Success {
		panic(checkErr)
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		logger.Error("AsyncSend Marshal", httpMethod, url, request, err)
		panic(err)
	}
	wantReq, err := http.NewRequest(strings.ToUpper(httpMethod), url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("AsyncSend NewRequest", httpMethod, url, jsonBody, err)
		panic(err)
	}
	wantReq.Header.Set("Accept", "application/json")
	wantReq.Header.Set("Content-Type", "application/json")

	go func() {
		// TODO: HTTP Request
	}()
}

func (w *WeCrossRPCService) checkRequest(request *methods.Request) *errors.Error {
	if request.Version == "" {
		return &errors.Error{Code: errors.RpcError, Detail: "Request version is empty"}
	} else {
		return &errors.Error{Code: errors.Success}
	}
}
