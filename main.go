package main

import (
	"github.com/gin-gonic/gin"
	"secureProxy/middleware"
	"secureProxy/proxy"
	"secureProxy/services"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/static", "./static")
	valkeyClient, _ := services.CreateClient()
	valkeyServ := services.NewValkeyService(valkeyClient)
	r.Use(middleware.ProxyMiddleware)

	r.GET("/auth", func(c *gin.Context) {
		redirectUrl := c.Query("redirectUrl")
		c.HTML(200, "login.html", gin.H{
			"RedirectUrl": redirectUrl,
		})
	})
	r.POST("/auth", func(c *gin.Context) {
		proxy.HandleAuthDomain(c, valkeyServ)
	})

	r.NoRoute(middleware.ProxyHandler)
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
