package fancy

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/** 抽取原先map为router */
type router struct {
	/** key 以及对应的方法 */
	handlers map[string]HandlerFunc

	/** 路由树 */
	roots    map[string]*node
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}


/** 解析路由将路由转化为part /api/v1/getUserId => ["api", "v1", "getUserId"] */
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


/** 添加路由 */
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {

	/** 解析路径 */
	parts := parsePattern(pattern)
	log.Printf("Route %4s - %s", method, pattern)

	key := method + "-" + pattern
	fmt.Println("key: ", key)

	/** 判断该方法数是否存在, 不存在则新建一个树 */
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	/** 插入前缀树 */
	r.roots[method].insert(pattern, parts, 0)

	/** 路由 key 对应的handler */
	r.handlers[key] = handler
}


func (r *router) getRoute(method string, path string) (*node, map[string]string) {

	/** 解析要匹配的path */
	searchParts := parsePattern(path)

	/** 存储动态路由的参数 */
	params := make(map[string]string)

	/** 方法数不存在则返回 */
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	/** 前缀树对比 */
	n := root.search(searchParts, 0)

	/** 处理匹配结果 */
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

/** 根据 context 处理 http */
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern

		/** r.handlers[key] 等同 HandlerFunc */
		c.handlers = append(c.handlers, r.handlers[key])

	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}