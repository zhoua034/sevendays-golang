package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
// newRouter create Router
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析pattern /p/go/doc 转为parts["p","go","doc"]
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute add route to r.roots and r.handlers
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// insert节点
	r.roots[method].insert(pattern, parts, 0)
	// save HandlerFunc
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// getRoute get route
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//解析请求的path
	searchParts := parsePattern(path)
	params := make(map[string]string)
	// 查找method对应的路由树根节点
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	// 根据root查找pattern
	n := root.search(searchParts, 0)
	if n != nil {
		//解析注册路由时的pattern
		parts := parsePattern(n.pattern)
		//封装params参数
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}
		return n, params
	}
	return nil, nil
}

// handle to handler context
func (r *router) handle(c *Context) {
	//解析请求，得到路由树的叶子节点，和请求参数params
	patternNode, parms := r.getRoute(c.Method, c.Path)
	if patternNode != nil {
		c.Params = parms
		key := c.Method + "-" + patternNode.pattern
		////从handlers中拿出注册时的HandlerFunc处理c
		//r.handlers[key](c)
		// 更新:
		//把当前请求的handler绑定到当前context的handlers中（也就是绑定在中间件之后）
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
