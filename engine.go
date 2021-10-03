package ws

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type HandlerFunc func(c *Context)
type MiddlewareList []HandlerFunc

type Engine struct {
	Manager   *Manager
	LocalAddr string
	RouteGroup
	ServerPrefix  string
	AccountPrefix string
	routeRule     map[string]*routeRule
	WsUserAddr
	ServerList *syncMap
}

func (e *Engine) GetEngine() *Engine {
	return e
}

func (e *Engine) SetNewManager() *Engine {
	e.Manager = newManager()
	return e
}

func (e *Engine) SetWsUserAddrFunc(w WsUserAddr) *Engine {
	e.WsUserAddr = w
	return e
}

func (e *Engine) SetLocalAddr(addr string) *Engine {
	e.LocalAddr = addr
	return e
}

func (e *Engine) SetServerPrefix(val string) *Engine {
	e.ServerPrefix = val
	return e
}

func (e *Engine) SetAccountPrefix(val string) *Engine {
	e.AccountPrefix = val
	return e
}

func (e *Engine) GetLocalAddr() string {
	return e.LocalAddr
}

func (e *Engine) StartManager() *Engine {
	e.Manager.start()
	return e
}

func (e *Engine) CheckSetAddr() bool {
	return e.WsUserAddr != nil
}

func (e *Engine) Run(w http.ResponseWriter, r *http.Request, uid UID) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Upgrade Error: %s", err.Error())
		return
	}

	socket.SetReadLimit(8192)

	if e.CheckSetAddr() {
		e.AddrExistWhenLoginHook(e, uid)
	} else {
		client, _ := e.Manager.GetClientByUid(uid)
		if client != nil {
			e.Manager.DisConnect <- client
		}
	}

	client := &Client{Socket: socket, Uid: uid, engine: e, GroupList: map[string]bool{}, SystemList: map[string]bool{}}
	e.Manager.addClient(client)

	if e.CheckSetAddr() {
		e.SetWsUserAddr(uid, e.LocalAddr)
	}

	client.read()
	e.Manager.Connect <- client
}

func (e *Engine) Use(m ...HandlerFunc) {
	e.middlewares = append(e.middlewares, m...)
}

func New() *Engine {
	engine := &Engine{}
	engine.engine = engine
	engine.ServerList = new(syncMap)
	engine.middlewares = make(MiddlewareList, 0)
	engine.routeRule = make(map[string]*routeRule)
	engine.SetNewManager().StartManager()

	return engine
}

type WsUserAddr interface {
	GetWsUserAddr(UID) string
	SetWsUserAddr(UID, string)
	DelWsUserAddr(UID)
	AddrExistWhenLoginHook(*Engine, UID)
}

// 添加客户端
func (e *Engine) AddServer(key, addr string) {
	_, loaded := e.ServerList.LoadOrStore(key, addr)
	if !loaded {
		e.ServerList.count++
	}
}

// 获取所有的客户端
func (e *Engine) GetServerList() []string {
	list := []string{}
	e.ServerList.Range(func(key, value interface{}) bool {
		list = append(list, value.(string))
		return true
	})
	return list
}

// 获取所有的客户端数量
func (e *Engine) GetServerCount() int {
	return e.ServerList.count
}

// 删除客户端
func (e *Engine) DelServer(key string) {
	_, loaded := e.ServerList.LoadAndDelete(key)
	if loaded {
		e.ServerList.count--
	}
}
