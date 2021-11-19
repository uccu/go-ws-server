package ws

import (
	"encoding/json"
	"sync"

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
	mu         sync.Mutex
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
					logrus.Warnf("close disConnect, clientId: %s", c.ClientId)
					return
				} else if messageType != websocket.PingMessage {
					logrus.Warnf("other disConnect, clientId: %s", c.ClientId)
					return
				}
			}
			if messageType == websocket.TextMessage {

				ackData := &Ack{}
				err := json.Unmarshal(message, ackData)
				if err != nil {
					logrus.Warnf("unmarshal err, clientId: %s,err: %s", c.ClientId, err.Error())
					continue
				}

				ack := ackData.Ack
				rule, ok := c.engine.routerRule[ack]
				if !ok {
					logrus.Warnf("no router, clientId: %s", c.ClientId)
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
	MessageId string      `json:"messageId"` // 消息ID
	Code      int         `json:"code"`      // 错误码
	Msg       string      `json:"msg"`       // 错误消息
	Data      interface{} `json:"data"`      // 数据内容
	Ack       string      `json:"ack"`       // 请求路由
}

func (c *Client) SendMessage(ack string, data interface{}) {

	if data == nil {
		data = struct{}{}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.Socket.WriteJSON(&RetData{
		MessageId: GetRandomStri(10),
		Ack:       ack,
		Data:      data,
	})
	if err != nil {
		logrus.Warnf("send message err,clientId: %s,err:%s", c.ClientId, err.Error())
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

	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.Socket.WriteJSON(&RetData{
		MessageId: GetRandomStri(10),
		Ack:       ack,
		Code:      code,
		Msg:       msg,
		Data:      data,
	})
	if err != nil {
		logrus.Warnf("send message err,clientId: %s,err:%s", c.ClientId, err.Error())
		c.engine.Manager.DisConnect <- c
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
