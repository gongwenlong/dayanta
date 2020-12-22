# Fancy Web Framework
Fancy is a simple web framework written in Go (Golang)

## Feature Overview

- Build robust and scalable RESTful APIs
- Context
- Group APIs
- Template rendering
- Extensible middleware framework
- Define middleware at group or route level
- Data binding for JSON and form payload
- Simple Centralized HTTP error handling
- Built in simple logger system

### Guide
## Installation
As of version v0.1.0, Fancy is available as a Go module, example

```sh
go get github.com/gongwenlong/dayanta v0.1.0

```

## Example

```sh
package main

import (
	fancy "github.com/gongwenlong/dayanta"
	"net/http"
)

func main()  {
	engine := fancy.New()

	v1 := engine.Group("/v1")
	{
		v1.GET("/hello", func(c *fancy.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	engine.Run(":9996")
}

```

## Router
```sh

 静态路由
 http://localhost:9996/json?username=xxxxx
  
 engine.GET("/json", func(c *fancy.Context) {
	 c.JSON(http.StatusOK, fancy.H{
		 "username": c.Query("username"),
	 })
 })
  
 engine.POST("/json", func(c *fancy.Context) {
	 c.JSON(http.StatusOK, fancy.H{
		 "username": c.PostForm("username"),
	 })
 })
 
 动态路由
 engine.GET("/hello/:name/login", func(c *fancy.Context) {
	 //  http://localhost:9996/hello/xxx/login
	 c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
 })

```

## Group APIS
```sh

 http://localhost:9996/v1/json?username=xxxxx
 
 v1 := engine.Group("/v1")
 {
	v1.GET("/json", func(c *fancy.Context) {
	   c.JSON(http.StatusOK, fancy.H{
		   "username": c.Query("username"),
	   })
   })
 }

```

## Middleware
```sh
func main()  {
 engine := fancy.New()
 engine.Use(Logger())
 
 v2 := r.Group("/v2")
 v2.Use(Logger()) // v2 group middleware
	{
		v2.GET("/hello", func(c *fancy.Context) {
			// http://localhost:9996/v2/hello
			c.String(http.StatusOK, c.Path)
		})
	}

 .....
}

func Logger() fancy.HandlerFunc {
	return func( c *fancy.Context) {
		t := time.Now()
		fmt.Println("time:" , t)

		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
```


## Panic
```sh
engine.GET("/panic", func(c *fancy.Context) {
	names := []string{"test"}
	c.String(http.StatusOK, names[100])
})

```

## Template Rendering 
```sh

http://localhost:9996/students

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main()  {
 engine := fancy.New()
 engine.SetFuncMap(template.FuncMap{
 "FormatAsDate": FormatAsDate,
 })
 engine.LoadHTMLGlob("templates/*")
 
 stu1 := &student{Name: "long", Age: 20}
 stu2 := &student{Name: "Jack", Age: 22}
 
 engine.GET("/students", func(c *fancy.Context) {
	c.HTML(http.StatusOK, "arr.tmpl", fancy.H{
	 "title":  "fancy",
	 "stuArr": [2]*student{stu1, stu2},
 	})
 })
 engine.Run(":9996")
}

项目根目录新建 templates

<!-- templates/arr.tmpl -->
<html>
<body>
    <p>hello, {{.title}}</p>
    {{range $index, $ele := .stuArr }}
    <p>{{ $index }}: {{ $ele.Name }} is {{ $ele.Age }} years old</p>
    {{ end }}
</body>
</html>

```

## Static
```sh
 访问localhost:9996/assets/222.png，最终返回/Users/admin/dongfangmingzhu/static/222.png
 engine.Static("/assets", "/Users/admin/dongfangmingzhu/static")
```

## License

[MIT](https://github.com/gongwenlong/dayanta/blob/master/LICENSE)
