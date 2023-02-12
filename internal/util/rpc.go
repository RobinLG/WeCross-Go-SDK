package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/methods"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	config "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"
)

func CheckPath(path string) error {
	sp := strings.Split(path, ".")
	if res, err := regexp.MatchString("^[A-Za-z]*.[A-Za-z0-9_-]*.[A-Za-z0-9_-]*$", path); !res || err != nil || len(sp) != 3 {
		return errors.New(fmt.Sprintf("Invalid iPath: %s", path), errors.ResourceError)
	}
	templateUrl := config.TEMPLATE_URL + strings.ReplaceAll(path, ".", "/")
	if _, err := url.ParseRequestURI(templateUrl); err != nil {
		return errors.New(fmt.Sprintf("Invalid iPath: %s", path), errors.IllegalSymbol)
	}
	return nil
}

func RecoverError(callback *methods.Callback) {
	switch err := recover().(type) {
	case errors.Error:
		callback.CallOnFailed(&errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("SDKError happened in AsyncSend, errorMessage: %s", err.Detail)})
	case error:
		callback.CallOnFailed(&errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("LibError happened in AsyncSend, errorMessage: %s", err.Error())})
	default:
	}
}

func FormatUrlPrefix(urlPrefix string) (string, *errors.Error) {
	pattern := "^/([^-_])[\\w-]{1,18}$"
	prefix := urlPrefix
	if len(prefix) == 0 {
		return "", nil
	}
	// /something/ => /something
	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}
	// something => /something
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	// /something
	matched, err := regexp.MatchString(pattern, prefix)
	if err != nil {
		return "", &errors.Error{Code: errors.FieldMissing, Detail: fmt.Sprintf("MatchString lib error: %s", err.Error())}
	}
	if matched {
		return prefix, nil
	} else {
		return "", &errors.Error{Code: errors.FieldMissing, Detail: fmt.Sprintf("Error 'urlPrefix' in config, it should matches pattern:  %s", pattern)}
	}
}
