package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/util"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("wecross_rpc")

const httpClientTimeout = 100000 // ms

type WeCrossRPCService struct {
	server     string
	httpClient *http.Client
}

func (w *WeCrossRPCService) InitService() error {
	w.httpClient = &http.Client{Timeout: httpClientTimeout * time.Millisecond}
	connection, err := w.getConnection(config.APPLICATION_CONFIG_FILE)
	if err != nil {
		return err
	}
	logger.Infof("connection: %v", connection)
	w.server = connection.Server
	return nil
}

func (w *WeCrossRPCService) Send(request *methods.Request, responseType methods.Response) (methods.Response, error) {
	return nil, nil
}

func (w *WeCrossRPCService) AsyncSend(request *methods.Request, responseType methods.Response, callback *methods.Callback) {
	defer util.RecoverError(callback)
	url := util.PathToUrl(w.server, request.Path) + "/" + request.Method

	checkErr := w.checkRequest(request)
	if checkErr.Code != errors.Success {
		panic(checkErr)
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		logger.Error("AsyncSend Marshal", url, request, err)
		panic(err)
	}
	wantReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("AsyncSend NewRequest", url, jsonBody, err)
		panic(err)
	}
	wantReq.Header.Set("Accept", "application/json")
	wantReq.Header.Set("Content-Type", "application/json")

	go func() {

	}()
}

func (w *WeCrossRPCService) checkRequest(request *methods.Request) *errors.Error {
	if request.Version == "" {
		return &errors.Error{Code: errors.RpcError, Detail: "Request version is empty"}
	} else {
		return &errors.Error{Code: errors.Success}
	}
}

func (w *WeCrossRPCService) getConnection(config string) (*Connection, *errors.Error) {
	connection, err := util.GetConnection(config)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
