package util

import (
	"fmt"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	config "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/config"
	"net/url"
	"regexp"
	"strings"
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
