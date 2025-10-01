package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"response": "pong",
		})
	})

	return r
}

func main() {
	router := setupRouter()

	err := router.RunTLS(":8443", "certs/secure-proxy-server-cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}
}
