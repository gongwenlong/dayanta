package fancy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

/**
Context随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载。
因此，设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。路由的处理函数，
以及将要实现的中间件，参数都统一使用 Context 实例，Context 就像一次会话的百宝箱，可以找到任何东西
 */
type Context struct {

	/** golang http原生 */
	Writer http.ResponseWriter
	Req    *http.Request

	/** 路径以及方法 */
	Path   string
	Method string

	StatusCode int

	/** 路由解析出来的参数 */
	Params map[string]string


	/** 中间件 */
	handlers []HandlerFunc
	index    int

	engine *Engine
}

/** 获取路由解析出来的参数的值 */
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

/** 创建context */
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index: -1,
	}
}

/*** 实现 Next 方法 */
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

/** 以下是获取参数 */
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}
