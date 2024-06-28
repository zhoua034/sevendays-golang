package gee

import (
	"log"
	"time"
)

// Logger middleware
func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		c.Next()
		// Calculate resolution time
		log.Printf("logger : [%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// OnlyV2 for v2 group test
func OnlyV2() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		c.Next()
		// Calculate resolution time
		log.Printf("OnlyV2 ：[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
