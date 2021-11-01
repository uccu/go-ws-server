package ws

import (
	"net/http"
)

type HandlerFunc func(c *Context)
type MiddlewareList []HandlerFunc

type Engine struct {
	Manager *Manager
	routerGroup
	routerRule map[string]*routerRule
}

func (e *Engine) GetEngine() *Engine {
	return e
}

func (e *Engine) GetManager() *Manager {
	return e.Manager
}

func (e *Engine) Run(w http.ResponseWriter, r *http.Request) *Client {

	socket := getSocket(w, r)
	if socket == nil {
		return nil
	}

	client := newClient(socket, e)
	e.Manager.addClient(client)

	client.read()
	e.Manager.connect <- client
	return client
}

func New() *Engine {
	engine := &Engine{}
	engine.engine = engine
	engine.Manager = newManager()
	engine.middlewares = make(MiddlewareList, 0)
	engine.routerRule = make(routerRuleMap)
	engine.Manager.start()

	return engine
}
