package resource

import (
	"fmt"
	internalwecrosslog "github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/wecrosslog"
)

var logger = wecrosslog.Component("resource")

type ResourceFactory struct{}

func (ResourceFactory) Build(weCrossRPC rpc.WeCrossRPC, path string) (*Resource, error) {
	resource := &Resource{
		weCrossRPC: weCrossRPC,
		path:       path,
	}
	resource.logger = internalwecrosslog.NewPrefixLogger(logger, fmt.Sprintf("[resource %p]", resource))
	return resource, resource.Check()
}
