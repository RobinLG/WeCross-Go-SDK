package util

import (
	"fmt"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	"github.com/pelletier/go-toml"
)

func GetToml(fileName string) (*toml.Tree, *errors.Error) {
	file, err := toml.LoadFile(fileName)
	if err != nil {
		return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("Something wrong with parsing %s : ", err.Error())}
	}
	return file, nil
}
