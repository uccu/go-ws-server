package ws

import (
	"net/http"
)

type HandlerFunc func(c *Context)
type MiddlewareList []HandlerFunc

type Engine struct {
	Manager *Manager
	RouteGroup
	routeRule map[string]*routeRule
}

type Hook struct {
	EngineHook func(*Engine)
	ClientHook func(*Engine, *Client)
}

func (e *Engine) GetEngine() *Engine {
	return e
}

func (e *Engine) GetManager() *Manager {
	return e.Manager
}

func (e *Engine) Run(w http.ResponseWriter, r *http.Request, hooks ...Hook) {

	socket := getSocket(w, r)
	if socket == nil {
		return
	}

	if len(hooks) > 0 && hooks[0].EngineHook != nil {
		hooks[0].EngineHook(e)
	}
	client := newClient(socket, e)
	e.Manager.addClient(client)

	if len(hooks) > 0 && hooks[0].EngineHook != nil {
		hooks[0].ClientHook(e, client)
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
