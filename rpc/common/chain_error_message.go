package common

type ChainErrorMessage struct {
	path    string
	message string
}

func (c *ChainErrorMessage) GetPath() string {
	return c.path
}

func (c *ChainErrorMessage) SetPath(path string) {
	c.path = path
}

func (c *ChainErrorMessage) GetMessage() string {
	return c.message
}

func (c *ChainErrorMessage) SetMessage(message string) {
	c.message = message
}
