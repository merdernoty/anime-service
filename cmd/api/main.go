package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// })
	// r.GET("/user/:name", func(c *gin.Context) {
	// })
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}