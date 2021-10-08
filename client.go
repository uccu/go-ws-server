package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Addr       string
	ClientId   string
	UnionId    string
	IsDeleted  bool
	Socket     *websocket.Conn
	GroupList  map[string]bool
	SystemList map[string]bool
	engine     *Engine
}

type Ack struct {
	Ack string `json:"ack"`
}

func (c *Client) read() {
	go func() {
		for {
			messageType, message, err := c.Socket.ReadMessage()
			if err != nil {
				if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					c.engine.Manager.disConnect <- c
					logrus.Debug("close disConnect, clientId: %d", c.ClientId)
					return
				} else if messageType != websocket.PingMessage {
					logrus.Debug("other disConnect, clientId: %d", c.ClientId)
					return
				}
			}
			if messageType == websocket.TextMessage {
				logrus.Debug("messageType:%d,message:%s", messageType, message)

				ackData := &Ack{}
				err := json.Unmarshal(message, ackData)
				if err != nil {
					logrus.Debug("unmarshal err, clientId: %d,err: %s", c.ClientId, err.Error())
					continue
				}

				ack := ackData.Ack
				rule, ok := c.engine.routeRule[ack]
				if !ok {
					logrus.Debug("no router, clientId: %d", c.ClientId)
					continue
				}

				ctx := NewContext(ack, c, message, rule)
				ctx.Engine = c.engine
				if len(rule.middlewares) > 0 {
					rule.middlewares[0](ctx)
				} else {
					rule.handleFunc(ctx)
				}
			}
		}
	}()
}

type RetData struct {
	MessageId string      `json:"messageId"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Ack       string      `json:"ack"`
}

func (c *Client) SendMessage(ack string, data interface{}) {

	if data == nil {
		data = struct{}{}
	}

	err := c.Socket.WriteJSON(&RetData{
		MessageId: GetRandomStri(10),
		Ack:       ack,
		Data:      data,
	})
	if err != nil {
		logrus.Debug("send message err,clientId: %d,err:%s", c.ClientId, err.Error())
		c.engine.Manager.disConnect <- c
	}
}

func (c *Client) SendErr(ack string, code int, msg string, datas ...interface{}) {

	var data interface{}
	if len(datas) > 0 {
		data = datas[0]
	}

	if data == nil {
		data = struct{}{}
	}

	err := c.Socket.WriteJSON(&RetData{
		MessageId: GetRandomStri(10),
		Ack:       ack,
		Code:      code,
		Msg:       msg,
		Data:      data,
	})
	if err != nil {
		logrus.Debug("send message err,clientId: %d,err:%s", c.ClientId, err.Error())
		c.engine.Manager.disConnect <- c
	}
}

func newClient(socket *websocket.Conn, engine *Engine) *Client {
	clientId := GetRandomStri(32)
	return &Client{Socket: socket,
		ClientId:   clientId,
		engine:     engine,
		GroupList:  map[string]bool{},
		SystemList: map[string]bool{},
	}
}
