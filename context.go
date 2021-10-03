package ws

import (
	"encoding/json"
	"errors"

	"github.com/sirupsen/logrus"
)

var ErrBind = errors.New("bind error")

func NewContext(ack string, client *Client, message []byte, rule *routeRule) *Context {
	context := &Context{}
	context.Ack = ack
	context.Client = client
	context.rule = rule
	context.CtxNo = GetRandomStri(6)
	context.RawMessage = message
	return context
}

type Context struct {
	CtxNo      string
	rule       *routeRule
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
	} else {
		c.rule.handleFunc(c)
	}
}

func (c *Context) ShouldBind(i interface{}) error {
	err := json.Unmarshal(c.RawMessage, i)
	if err != nil {
		logrus.Debug("uid: %d,err:%s", c.Client.Uid, err.Error())
		return ErrBind
	}
	return nil
}

func (c *Context) GetLocalAddr() string {
	return c.Engine.LocalAddr
}
func (c *Context) GetEngine() *Engine {
	return c.Engine
}
