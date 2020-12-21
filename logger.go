package fancy

import (
	"fmt"
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		fmt.Println("time:" , t)

		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}