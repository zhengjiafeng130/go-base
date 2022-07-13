package base

import "log"

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // nesting
	engine      *Engine 	  // all router group shares an engine instance
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix:      prefix,
		parent:      g,
		engine:      g.engine,
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute(HttpMethodGET, pattern, handler)
}

// POST defines the method to add POST request
func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute(HttpMethodPOST, pattern, handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}