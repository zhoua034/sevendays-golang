package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {

	r := gee.Default()
	r.Static("/assets", "./static")

	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.html", nil)
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})
	v1 := r.Group("/v1")
	{
		// expect /hello?name=geektutu
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s ,you are at%s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(gee.OnlyV2())
	{
		v2.POST("/hello/:name", func(c *gee.Context) {
			// expect /hello/xxxx
			c.String(http.StatusOK, "hello %s ,you are at%s\n", c.Param("name"), c.Path)
		})
		//v2.POST("/login", func(c *gee.Context) {
		//	c.JSON(http.StatusOK, gee.H{
		//		"username": c.PostForm("username"),
		//		"password": c.PostForm("password"),
		//	})
		//})
	}

	fmt.Println("Gee 启动成功！！！")
	r.Run(":8085")
}
