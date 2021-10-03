package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type UID = int64
type Client struct {
	Addr       string
	Uid        UID
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
					c.engine.Manager.DisConnect <- c
					logrus.Debug("close disConnect, uid: %d", c.Uid)
					return
				} else if messageType != websocket.PingMessage {
					logrus.Debug("other disConnect, uid: %d", c.Uid)
					return
				}
			}
			if messageType == websocket.TextMessage {
				logrus.Debug("messageType:%d,message:%s", messageType, message)

				ackData := &Ack{}
				err := json.Unmarshal(message, ackData)
				if err != nil {
					logrus.Debug("unmarshal err, uid: %d,err: %s", c.Uid, err.Error())
					continue
				}

				ack := ackData.Ack
				rule, ok := c.engine.routeRule[ack]
				if !ok {
					logrus.Debug("no router, uid: %d", c.Uid)
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
		logrus.Debug("send message err,uid: %d,err:%s", c.Uid, err.Error())
		c.engine.Manager.DisConnect <- c
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
		logrus.Debug("send message err,uid: %d,err:%s", c.Uid, err.Error())
		c.engine.Manager.DisConnect <- c
	}
}
