package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secureProxy/appConfig"
	"secureProxy/proxy"
	"secureProxy/services"
	"strings"
)

var config = appConfig.CreateConfig()
var valkeyClient, _ = services.CreateClient()
var valkeyServ = services.NewValkeyService(valkeyClient)

func ProxyMiddleware(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	if host == config.AuthDomain {
		log.Println("AuthDomain request was accepted")
		proxy.HandleAuthDomain(c, valkeyServ)
		c.JSON(http.StatusOK, gin.H{
			"loginStatus": "success",
		})
		//c.Next()
		return
	}

	if !checkAuthentication(c) {
		return
	}

	proxyToUpstream(c, host)
}

func checkAuthentication(c *gin.Context) bool {
	sessionCookie, err := c.Cookie("SECURE_PROXY_SESSION")
	if err != nil {
		redirectToAuth(c)
		return false
	}

	email, err := valkeyServ.Get(c, "session:"+sessionCookie)
	if err != nil {
		redirectToAuth(c)
		return false
	}

	valkeyServ.Expire(c, "session:"+sessionCookie, 1800) // 30 минут

	c.Set("authenticated_user", email)

	return true
}

func proxyToUpstream(c *gin.Context, host string) {
	for i := range config.Upstreams {
		upstream := &config.Upstreams[i]
		if upstream.Host == host {
			proxyRequest(upstream, c)
			return
		}
	}

	c.AbortWithError(http.StatusNotFound, errors.New("upstream not found"))
}

func redirectToAuth(c *gin.Context) {
	authenticatedRedirectUrl := "https://" + c.Request.Host + c.Request.RequestURI
	authUrl := "https://" + config.AuthDomain + "/auth?redirectUrl=" + url.QueryEscape(authenticatedRedirectUrl)
	c.Redirect(http.StatusFound, authUrl)
}

func proxyRequest(upstream *appConfig.Upstream, c *gin.Context) {
	upstreamUrl, err := url.Parse(upstream.Destination + c.Request.RequestURI)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(upstreamUrl)
	proxy.Director = func(req *http.Request) {
		req.URL = upstreamUrl
		if user, exists := c.Get("authenticated_user"); exists {
			req.Header.Set("X-Forwarded-User", user.(string))
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
