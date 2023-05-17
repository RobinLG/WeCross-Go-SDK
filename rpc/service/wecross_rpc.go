package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/authentication"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/util"
	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	req "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods/request"
	res "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods/response"
)

const (
	prefix            = "[rpc-service] "
	httpClientTimeout = 10 * time.Second
)

type WeCrossRPCService struct {
	server     string
	httpClient *http.Client
	urlPrefix  string
	logger     *internalwecrosslog.PrefixLogger
}

type completableFuture struct {
	response interface{}
	error    *errors.Error
}

func (w *WeCrossRPCService) InitService(l internalwecrosslog.DepthLoggerV1) *errors.Error {
	w.logger = internalwecrosslog.NewPrefixLogger(l, prefix)
	connection, err := w.getConnection(config.APPLICATION_CONFIG_FILE)
	if err != nil {
		return err
	}
	w.logger.Infof("connection: %v", connection)
	if connection.SslSwitch == config.SSL_OFF {
		w.server = "http://" + connection.Server
	} else {
		w.server = "https://" + connection.Server
	}
	if len(connection.UrlPrefix) > 0 {
		w.urlPrefix = connection.UrlPrefix
	}
	if w.httpClient, err = w.getHttpClient(connection); err != nil {
		return err
	}
	return nil
}

func (w *WeCrossRPCService) Send(httpMethod string, uri string, request *methods.Request, response methods.Response) (methods.Response, error) {
	w.checkRequest(request)
	completableFutureCh := make(chan completableFuture)

	callback := methods.CallbackFactory{}.Build()
	callback.OnSuccess = func(response methods.Response) {
		completableFutureCh <- completableFuture{response: response, error: nil}
	}
	callback.OnFailed = func(err *errors.Error) {
		completableFutureCh <- completableFuture{response: nil, error: err}
	}
	w.AsyncSend(httpMethod, uri, request, response, callback)

	var cf completableFuture
	select {
	case <-time.After(20 * time.Second):
		break
	case cf = <-completableFutureCh:
		break
	}

	w.logger.Debugf("response: %v", response)

	if cf.error != nil {
		return methods.Response{}, cf.error
	}

	if uaRes, ok := cf.response.(res.UAResponse); ok {
		if uaReq, ok1 := request.GetData().(req.UARequest); ok1 {
			w.GetUAResponseInfo(uri, uaReq, uaRes)
		} else if uaReq, ok1 = request.GetExt().(req.UARequest); ok1 {
			w.GetUAResponseInfo(uri, uaReq, uaRes)
		}
	}

	return methods.Response{}, nil
}

func (w *WeCrossRPCService) GetUAResponseInfo(uri string, uaRequest req.UARequest, response res.UAResponse) errors.Error {
	query := strings.Split(uri, "/")[2]
	if query == "login" {
		credential := response.GetUAReceipt().GetCredential()

		w.logger.Infof("CurrentUse: %s", uaRequest.GetUsername())
		if credential == "" {
			w.logger.Errorf("Login fail, credential in UAResponse is null")
			return errors.Error{Code: errors.RpcError, Detail: "Login fail, credential in UAResponse is null!"}
		}
		authentication.SetCurrentUser(uaRequest.GetUsername(), credential)
	}
	if query == "logout" {
		w.logger.Infof("CurrentUser: %s logout.", authentication.GetCurrentUser())
		authentication.ClearCurrentUser()
	}
	return errors.Error{}
}

func (w *WeCrossRPCService) AsyncSend(httpMethod, uri string, request *methods.Request, response methods.Response, callback *methods.Callback) {
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
		w.logger.Errorf("AsyncSend Marshal: %s, %v, %s", url, request, err.Error())
		panic(err)
	}
	httpRequest, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		w.logger.Errorf("AsyncSend NewRequest: %s, %v, %s", url, jsonBody, err)
		panic(err)
	}
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	go func() {
		httpResponse, errDo := w.httpClient.Do(httpRequest)
		if errDo != nil {
			panic(&errors.Error{Code: errors.InternalError, Detail: "handle response failed: " + errDo.Error()})
		}

		if httpResponse.StatusCode == 401 {
			callback.CallOnFailed(
				&errors.Error{
					Code: errors.LackAuthentication,
					Detail: "HTTP status code: 401-Unauthorized, have you logged in?\n" +
						"If you have logged-in already, maybe you should re-login " +
						"because your account login status has expired.",
				})
			return
		}
		if httpResponse.StatusCode == 404 {
			callback.CallOnFailed(
				&errors.Error{
					Code: errors.LackAuthentication,
					Detail: "HTTP status code: 404 Not Found\n" +
						"Maybe your request's resource path is wrong.",
				})
			return
		}
		if httpResponse.StatusCode != 200 {
			callback.CallOnFailed(&errors.Error{
				Code:   errors.RpcError,
				Detail: fmt.Sprintf("HTTP response status: %d message: %s", httpResponse.StatusCode, httpResponse.Status),
			})
			return
		} else {
			buf := &bytes.Buffer{}
			buf.ReadFrom(httpResponse.Body)

			// TODO *response according type
			errJson := json.Unmarshal(buf.Bytes(), response)
			if errJson != nil {
				panic(&errors.Error{Code: errors.InternalError, Detail: "HTTP response status is 200, but unmarshal error. Detail: " + errJson.Error()})
			}

			callback.CallOnSuccess(response)
		}
	}()
}

func (w *WeCrossRPCService) checkRequest(request *methods.Request) *errors.Error {
	if request.GetVersion() == "" {
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
	connection.CaCert = file.GetString("connection.caCert")
	connection.SslKey = file.GetString("connection.sslKey")
	connection.SslCert = file.GetString("connection.sslCert")
	connection.SslSwitch = file.GetInt64("connection.sslSwitch")
	if connection.UrlPrefix, err = util.FormatUrlPrefix(file.GetString("connection.urlPrefix")); err != nil {
		return nil, err
	}
	return connection, nil
}

func (w *WeCrossRPCService) getServer(toml *util.WeCrossToml) (string, *errors.Error) {
	server := toml.GetString("connection.server")
	if len(server) == 0 {
		return "", &errors.Error{Code: errors.FieldMissing, Detail: "Something wrong with parsing [connection.server], please check configuration"}
	}
	return server, nil
}

func (w *WeCrossRPCService) getHttpClient(connection *Connection) (*http.Client, *errors.Error) {
	transport := &http.Transport{
		DisableKeepAlives:     false,
		TLSHandshakeTimeout:   httpClientTimeout,
		IdleConnTimeout:       httpClientTimeout,
		ResponseHeaderTimeout: httpClientTimeout,
		ExpectContinueTimeout: httpClientTimeout,
	}
	if connection.SslSwitch != config.SSL_OFF {
		caCert, err := ioutil.ReadFile(connection.CaCert)
		if err != nil {
			w.logger.Errorf("Init http client error: %s", err.Error())
			return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("Init http client error: %s", err.Error())}
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		cert, err := tls.LoadX509KeyPair(connection.SslCert, connection.SslKey)
		if err != nil {
			w.logger.Errorf("Init http client error: %s", err.Error())
			return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("Init http client error: %s", err.Error())}
		}
		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}

	clint := &http.Client{
		Transport: transport,
		Timeout:   httpClientTimeout,
	}

	return clint, nil
}
