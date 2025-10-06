package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"secureProxy/services"
)

func HandleAuthDomain(c *gin.Context, valkeyServ *services.ValkeyService) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "login.html", gin.H{})
		return
	}
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error": "Email is required",
		})
		return
	}
	if password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error": "Password is required",
		})
		return
	}
	secret := services.GenerateTOTP(c, email)
	sessionCookie := SetProxyCookie(c)
	valkeyServ.Set(c, "session:"+sessionCookie, secret)
	valkeyServ.Expire(c, email, 1800)
}

func SetProxyCookie(c *gin.Context) string {
	sessionId := uuid.New().String()
	c.SetCookie("SECURE_PROXY_SESSION", sessionId, 60*60, "/", ".secure-proxy.lan", true, true)
	return sessionId
}
