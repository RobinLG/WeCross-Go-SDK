package util

import (
	"fmt"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/pelletier/go-toml"
)

func GetToml(fileName string) (*WeCrossToml, *errors.Error) {
	file, err := toml.LoadFile(fileName)
	if err != nil {
		return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("Something wrong with parsing %s : ", err.Error())}
	}
	return &WeCrossToml{file}, nil
}

type WeCrossToml struct {
	*toml.Tree
}

func (t *WeCrossToml) GetString(key string) string {
	val := t.Get(key)
	if val == nil {
		return ""
	} else {
		if v, ok := val.(string); ok {
			return v
		} else {
			return ""
		}
	}
}

func (t *WeCrossToml) GetInt64(key string) int64 {
	val := t.Get(key)
	if val == nil {
		return 0
	} else {
		if v, ok := val.(int64); ok {
			return v
		} else {
			return 0
		}
	}
}
