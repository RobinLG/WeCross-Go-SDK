package methods

import (
	"sync/atomic"
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/errors"
)

const (
	CALLBACK_TIMEOUT = 30000
)

func newCallback() *Callback {
	c := &Callback{
		timeout:  CALLBACK_TIMEOUT,
		isFinish: atomic.Value{},
	}
	c.isFinish.Store(false)

	c.timeoutWorker = func() {
		timer := time.NewTimer(CALLBACK_TIMEOUT * time.Millisecond)
		<-timer.C
		if c.isFinish.Swap(true) == false { // timeout
			c.err = errors.Error{
				Code:   errors.RemoteCallError,
				Detail: "Timeout",
			}
		}
	}

	return c
}

type Callback struct {
	timeout       int
	timeoutWorker func()
	isFinish      atomic.Value
	err           errors.Error
	callbackWorker
}

type callbackWorker interface {
	OnSuccess(func(Response))
	OnFailed(func(errors.Error))
}
