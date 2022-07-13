package base

import (
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots: map[string]*node{},
		handlers: map[string]HandlerFunc{},
	}
}

// * is allowed
func parsePattern(pattern string) []string {
	partSplits := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range partSplits {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//log.Printf("Route %4s - %s", method, pattern)

	key := method + "-" + pattern

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parsePattern(pattern), 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string) // 模糊匹配项

	var root *node
	var ok bool
	if root, ok = r.roots[method]; !ok {
		return nil, nil
	}

	if _node := root.search(searchParts, 0); _node != nil {
		parts := parsePattern(_node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return _node, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	//key := c.Method + "-" + c.Path
	if _node, params := r.getRoute(c.Method, c.Path); _node != nil {
		c.Params = params
		key := c.Method + "-" + _node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}