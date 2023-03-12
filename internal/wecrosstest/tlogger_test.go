package wecrosstest

import (
	"testing"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrosslog"
)

type s struct {
	Tester
}

func Test(t *testing.T) {
	RunSubTests(t, s{})
}

func (s) TestInfo(t *testing.T) {
	wecrosslog.Logger.Info("Info", "message.")
}

func (s) TestInfoln(t *testing.T) {
	wecrosslog.Logger.Infoln("Info", "message.")
}

func (s) TestInfof(t *testing.T) {
	wecrosslog.Logger.Infof("%v %v.", "Info", "message")
}

func (s) TestInfoDepth(t *testing.T) {
	wecrosslog.InfoDepth(0, "Info", "depth", "message.")
}

func (s) TestWarning(t *testing.T) {
	wecrosslog.Logger.Warning("Warning", "message.")
}

func (s) TestWarningln(t *testing.T) {
	wecrosslog.Logger.Warningln("Warning", "message.")
}

func (s) TestWarningf(t *testing.T) {
	wecrosslog.Logger.Warningf("%v %v.", "Warning", "message")
}

func (s) TestWarningDepth(t *testing.T) {
	wecrosslog.WarningDepth(0, "Warning", "depth", "message.")
}

func (s) TestError(t *testing.T) {
	const numErrors = 10
	TLogger.ExpectError("Expected error")
	TLogger.ExpectError("Expected ln error")
	TLogger.ExpectError("Expected formatted error")
	TLogger.ExpectErrorN("Expected repeated error", numErrors)
	wecrosslog.Logger.Error("Expected", "error")
	wecrosslog.Logger.Errorln("Expected", "ln", "error")
	wecrosslog.Logger.Errorf("%v %v %v", "Expected", "formatted", "error")
	for i := 0; i < numErrors; i++ {
		wecrosslog.Logger.Error("Expected repeated error")
	}
}
