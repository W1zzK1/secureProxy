package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secureProxy/middleware"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ProxyMiddleware)

	getGroup := r.Group("/")
	{
		getGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"response": "pong",
			})
		})
		getGroup.GET("/set_cookie", func(c *gin.Context) {
			c.SetCookie("SECURE_PROXY_SESSION", "E364EEAE-8F50-4B6E-BB9B-E7F56A27160C", 60*60*24, "/", ".secure-proxy.lan", true, true)
			c.JSON(200, gin.H{
				"cookie_val": "success",
			})
		})
		getGroup.GET("/proxy", func(c *gin.Context) {
			parsedUrl, err := url.Parse("http://localhost:8181/") // bank-vsm-restaurant

			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			proxy := httputil.NewSingleHostReverseProxy(parsedUrl)
			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = parsedUrl.Scheme
				req.URL.Host = parsedUrl.Host
				req.URL.Path = parsedUrl.Path + "payment/fail"
				req.Header.Set("X-Forwarded_User", "w1zzk1")
			}
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	return r
}

func main() {
	router := setupRouter()

	err := router.RunTLS(":8443", "certs/secure-proxy-server-cert.pem", "certs/key.pem")
	if err != nil {
		panic(err)
	}
}
