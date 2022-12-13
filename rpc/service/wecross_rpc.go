package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"
	"github.com/pelletier/go-toml"

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
	urlPrefix  string
}

func (w *WeCrossRPCService) InitService() *errors.Error {
	connection, err := w.getConnection(config.APPLICATION_CONFIG_FILE)
	if err != nil {
		return err
	}
	logger.Infof("connection: %v", connection)
	if connection.SslSwitch == config.SSL_OFF {
		w.server = "http://" + connection.Server
	} else {
		w.server = "https://" + connection.Server
	}
	if len(connection.UrlPrefix) > 0 {
		w.urlPrefix = connection.UrlPrefix
	}
	if w.httpClient, err = w.getHttpAsyncClient(connection); err != nil {
		return err
	}
	return nil
}

func (w *WeCrossRPCService) Send(request *methods.Request, responseType methods.Response) (methods.Response, error) {
	return nil, nil
}

func (w *WeCrossRPCService) AsyncSend(httpMethod, uri string, request *methods.Request, responseType methods.Response, callback *methods.Callback) {
	defer util.RecoverError(callback)
	url := ""
	if len(w.urlPrefix) > 0 {
		url = w.server + w.urlPrefix + uri
	} else {
		url = w.server + uri
	}

	checkErr := w.checkRequest(request)
	if checkErr.Code != errors.Success {
		panic(checkErr)
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		logger.Error("AsyncSend Marshal", url, request, err)
		panic(err)
	}
	wantReq, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBody))
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
	file, err := util.GetToml(config)
	if err != nil {
		return nil, err
	}
	connection := &Connection{}
	if connection.Server, err = w.getServer(file); err != nil {
		return nil, err
	}
	connection.CaCert = fmt.Sprintf("%s", file.Get("connection.caCert"))
	connection.SslKey = fmt.Sprintf("%s", file.Get("connection.sslKey"))
	connection.SslCert = fmt.Sprintf("%s", file.Get("connection.sslCert"))
	connection.SslSwitch = file.Get("connection.sslSwitch").(int)
	if connection.UrlPrefix, err = util.FormatUrlPrefix(fmt.Sprintf("%s", file.Get("connection.urlPrefix"))); err != nil {
		return nil, err
	}
	return connection, nil
}

func (w *WeCrossRPCService) getServer(toml *toml.Tree) (string, *errors.Error) {
	server := fmt.Sprintf("%s", toml.Get("connection.server"))
	if len(server) == 0 {
		return "", &errors.Error{Code: errors.FieldMissing, Detail: "Something wrong with parsing [connection.server], please check configuration"}
	}
	return server, nil
}

func (w *WeCrossRPCService) getHttpAsyncClient(connection *Connection) (*http.Client, *errors.Error) {
	return nil, nil
}
