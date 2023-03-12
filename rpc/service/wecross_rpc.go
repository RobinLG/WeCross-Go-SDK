package service

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/util"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("wecross_rpc")

const httpClientTimeout = 10 * time.Second

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
	if w.httpClient, err = w.getHttpClient(connection); err != nil {
		return err
	}
	return nil
}

func (w *WeCrossRPCService) Send(httpMethod string, uri string, request *methods.Request, responseType methods.Response) (methods.Response, error) {
	return nil, nil
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
		logger.Error("AsyncSend Marshal", url, request, err)
		panic(err)
	}
	httpRequest, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("AsyncSend NewRequest", url, jsonBody, err)
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

			errJson := json.Unmarshal(buf.Bytes(), response)
			if errJson != nil {
				panic(&errors.Error{Code: errors.InternalError, Detail: "HTTP response status is 200, but unmarshal error. Detail: " + errJson.Error()})
			}

			callback.CallOnSuccess(response)
		}
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
			logger.Error("Init http client error: ", err)
			return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("Init http client error: %s", err.Error())}
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		cert, err := tls.LoadX509KeyPair(connection.SslCert, connection.SslKey)
		if err != nil {
			logger.Error("Init http client error: ", err)
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
