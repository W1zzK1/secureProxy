package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"secureProxy/appConfig"
)

var config = appConfig.CreateConfig()

func ProxyMiddleware(c *gin.Context) {
	host := c.GetHeader("Host")
	if host == config.AuthDomain {
		c.Next()
	}
	_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Unauthorized"))
}
