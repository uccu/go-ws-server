package ws

import (
	"strings"
)

type Router interface {
	Group(path string, m ...HandlerFunc) *routerGroup
	GET(path string, handleFunc HandlerFunc)
	Use(m ...HandlerFunc)
}

type routerGroup struct {
	basePath    string
	engine      *Engine
	upperGroup  *routerGroup
	middlewares MiddlewareList
}

type routerRule struct {
	handleFunc  HandlerFunc
	middlewares MiddlewareList
}

type routerRuleMap map[string]*routerRule

func (e *routerGroup) Use(m ...HandlerFunc) {
	e.middlewares = append(e.middlewares, m...)
}

func (e *routerGroup) Group(path string, m ...HandlerFunc) *routerGroup {

	path = strings.ReplaceAll(path, "/", ".")
	path = strings.Trim(path, ".")

	group := &routerGroup{
		basePath:    path,
		engine:      e.engine,
		upperGroup:  e,
		middlewares: e.middlewares,
	}
	group.middlewares = append(group.middlewares, m...)
	return group
}

func (e *routerGroup) GET(path string, handleFunc HandlerFunc) {
	path = strings.ReplaceAll(path, "/", ".")
	path = strings.Trim(path, ".")
	group := e
	for group != nil {
		if group.basePath != "" {
			path = group.basePath + "." + path
		}
		group = group.upperGroup
	}
	path = strings.Trim(path, ".")
	e.engine.routerRule[path] = &routerRule{handleFunc: handleFunc, middlewares: e.middlewares}
}
