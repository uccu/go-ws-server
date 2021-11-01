package ws

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
)

var ErrBind = errors.New("bind error")

func NewContext(ack string, client *Client, message []byte, rule *routerRule) *Context {
	context := &Context{}
	context.Ack = ack
	context.Client = client
	context.rule = rule
	context.RequestId = GetRandomStri(6)
	context.RawMessage = message
	return context
}

type Context struct {
	RequestId  string
	rule       *routerRule
	Ack        string
	Engine     *Engine
	Client     *Client
	RawMessage []byte
	nm         int
}

func (c *Context) Next() {
	c.nm++
	if c.nm < len(c.rule.middlewares) {
		c.rule.middlewares[c.nm](c)
	} else if c.nm == len(c.rule.middlewares) {
		c.rule.handleFunc(c)
	}
}

func (c *Context) ShouldBind(i interface{}) error {
	err := json.Unmarshal(c.RawMessage, i)
	if err != nil {
		logrus.Warnf("clientId: %s,err:%s", c.Client.ClientId, err.Error())
		return ErrBind
	}
	return nil
}

func (c *Context) GetEngine() *Engine {
	return c.Engine
}
