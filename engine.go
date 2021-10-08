package ws

import (
	"net/http"
)

type HandlerFunc func(c *Context)
type MiddlewareList []HandlerFunc
type Hook func(*Engine)

type Engine struct {
	Manager *Manager
	RouteGroup
	routeRule map[string]*routeRule
}

func (e *Engine) GetEngine() *Engine {
	return e
}

func (e *Engine) GetManager() *Manager {
	return e.Manager
}

func (e *Engine) Run(w http.ResponseWriter, r *http.Request, hook ...Hook) {

	socket := getSocket(w, r)
	if socket == nil {
		return
	}

	if len(hook) > 0 {
		hook[0](e)
	}
	client := newClient(socket, e)
	e.Manager.addClient(client)

	if len(hook) > 1 {
		hook[1](e)
	}

	client.read()
	e.Manager.Connect <- client

}

func New() *Engine {
	engine := &Engine{}
	engine.engine = engine
	engine.Manager = newManager()
	engine.middlewares = make(MiddlewareList, 0)
	engine.routeRule = make(routeRuleMap, 0)
	engine.Manager.start()

	return engine
}
