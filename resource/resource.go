package resource

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
	internalutil "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/util"
	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc"
)

type Resource struct {
	weCrossRPC rpc.WeCrossRPC
	path       string
	logger     *internalwecrosslog.PrefixLogger
}

func (r *Resource) Check() error {
	if err := r.checkWeCrossRPC(); err != nil {
		return err
	}
	if err := r.checkIPath(); err != nil {
		return err
	}
	return nil
}

func (r *Resource) checkWeCrossRPC() error {
	if r.weCrossRPC == nil {
		return errors.New("WeCrossRPC not set", errors.ResourceError)
	}
	return nil
}

func (r *Resource) checkIPath() error {
	return internalutil.CheckPath(r.path)
}
