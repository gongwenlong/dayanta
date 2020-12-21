package fancy

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

/** 定义函数 */
type HandlerFunc func(*Context)

/*** Engine拥有RouterGroup所有的能力 */
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有的分组

	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render
}

/** 定义框架 */
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

type RouterGroup struct {
	prefix      string        // 分组前缀
	middlewares []HandlerFunc // 中间键
	parent      *RouterGroup  // 嵌套
	engine      *Engine       // engine 指针
}

/** 创建分组 */
/** 所有的group共享一个engine */
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

/** 添加路由 */
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

/** 添加Get方法 */
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

/** 添加Post方法 */
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}


// 创建静态 handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	fmt.Println(relativePath)
	// 绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")

		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 静态服务
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 添加路由
	group.GET(urlPattern, handler)
}

// html 模板
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

/**
Engine 实现 ServeHTTP
Engine 等于 Handler

func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
*/
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}



// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

/**

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

Engine 实现 ServeHTTP
Engine 等于 Handler
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	/** 从分组添加中间件到 middlewares */
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	/** 拿到请求的url */
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine

	/** 根据请求的url获取 */
	engine.router.handle(c)
}