package util

import (
	"fmt"
	"io/ioutil"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	rpc "github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service"
	"github.com/pelletier/go-toml"
)

func GetConnection(fileName string) (*rpc.Connection, *errors.Error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("io read wrong with parsing %s : ", err.Error())}
	}
	connection := &rpc.Connection{}
	err = toml.Unmarshal(file, connection)
	if err != nil {
		return nil, &errors.Error{Code: errors.InternalError, Detail: fmt.Sprintf("toml unmarshle wrong with parsing %s : ", err.Error())}
	}
	return connection, nil
}
