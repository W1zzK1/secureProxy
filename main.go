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

	r.GET("/set_cookie", func(c *gin.Context) {
		c.SetCookie("SECURE_PROXY_SESSION", "E364EEAE-8F50-4B6E-BB9B-E7F56A27160C", 60*60*24, "/", ".secure-proxy.lan", true, true)
		c.JSON(200, gin.H{
			"cookie_val": "success",
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
