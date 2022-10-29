package methods

type CallbackFactory struct{}

func (c CallbackFactory) Build() *Callback {
	return newCallback()
}
