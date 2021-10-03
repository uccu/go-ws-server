package ws

import (
	"strings"
)

type RouteGroup struct {
	basePath    string
	engine      *Engine
	upperGroup  *RouteGroup
	middlewares MiddlewareList
}

type routeRule struct {
	handleFunc  HandlerFunc
	middlewares MiddlewareList
}

func (e *RouteGroup) Group(path string, m ...HandlerFunc) *RouteGroup {

	path = strings.ReplaceAll(path, "/", ".")
	path = strings.Trim(path, ".")

	group := &RouteGroup{
		basePath:    path,
		engine:      e.engine,
		upperGroup:  e,
		middlewares: e.middlewares,
	}
	group.middlewares = append(group.middlewares, m...)
	return group
}

func (e *RouteGroup) GET(path string, handleFunc HandlerFunc) {
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
	e.engine.routeRule[path] = &routeRule{handleFunc: handleFunc, middlewares: e.middlewares}
}
