package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secureProxy/appConfig"
	"secureProxy/services"
	"strings"
)

var config = appConfig.CreateConfig()

var valkeyClient, _ = services.CreateClient()
var valkeyServ = services.NewValkeyService(valkeyClient)

func ProxyMiddleware(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	if host == config.AuthDomain || c.Request.URL.Path == "/auth" {
		log.Println("AuthDomain request was accepted")

		c.Next()
		return
	}

	if !checkAuthentication(c, valkeyServ) {
		c.Abort()
		return
	}

	c.Next()
}

func ProxyHandler(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	if host == config.AuthDomain {
		c.Next()
		return
	}

	ProxyToUpstream(c, host)
}

func checkAuthentication(c *gin.Context, valkeyServ *services.ValkeyService) bool {
	sessionCookie, err := c.Cookie("SECURE_PROXY_SESSION")
	if err != nil {
		redirectToAuth(c)
		return false
	}

	username, err := valkeyServ.Get(c, "session:"+sessionCookie)
	if err != nil {
		redirectToAuth(c)
		return false
	}

	valkeyServ.Expire(c, "session:"+sessionCookie, 1800) // 30 минут

	c.Set("authenticated_user", username)

	return true
}

func ProxyToUpstream(c *gin.Context, host string) {
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
	authUrl := "https://" + config.AuthDomain + ":8443/auth?redirectUrl=" + url.QueryEscape(authenticatedRedirectUrl)
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
		req.URL.Scheme = upstreamUrl.Scheme
		req.URL.Host = upstreamUrl.Host
		req.URL.Path = upstreamUrl.Path
		req.URL.RawQuery = c.Request.URL.RawQuery
		if user, exists := c.Get("authenticated_user"); exists {
			req.Header.Set("X-Forwarded-User", user.(string))
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
