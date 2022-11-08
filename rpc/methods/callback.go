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
		Timeout:  CALLBACK_TIMEOUT,
		IsFinish: atomic.Value{},
	}
	c.IsFinish.Store(false)

	c.Timer = time.AfterFunc(CALLBACK_TIMEOUT*time.Millisecond,
		func() {
			if c.IsFinish.Swap(true) == false { // timeout
				c.CallOnFailed(&errors.Error{
					Code:   errors.RemoteCallError,
					Detail: "Timeout",
				})
			}
		})

	return c
}

type Callback struct {
	Timeout   int
	Timer     *time.Timer
	IsFinish  atomic.Value
	OnFailed  func(*errors.Error)
	OnSuccess func(Response)
}

func (c *Callback) CallOnFailed(error *errors.Error) {
	if c.IsFinish.Swap(true) == false {
		c.Timer.Stop()
		c.OnFailed(error)
	}
}
