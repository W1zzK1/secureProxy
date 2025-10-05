package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secureProxy/appConfig"
	"secureProxy/valkeyService"
	"strings"
)

var config = appConfig.CreateConfig()
var valkeyClient, _ = valkeyService.CreateClient()
var valkeyServ = valkeyService.NewValkeyService(valkeyClient)

func ProxyMiddleware(c *gin.Context) {
	host := strings.Split(c.Request.Host, ":")[0]
	if host == config.AuthDomain {
		log.Println("request was accepted")
		c.SetCookie("username", "vr.gorbunov", 60*60*24, "/", ".secure-proxy.lan", true, true)
		cookie, _ := c.Cookie("username")
		c.JSON(http.StatusOK, gin.H{
			"cookieValue": cookie,
		})
		valkeyServ.Set(c, "username", "vr.gorbunov")
		c.Next()
		return
	}

	c.Abort()
	getUserName(c)

	for i := range config.Upstreams {
		upstream := &config.Upstreams[i]
		if upstream.Host == host {
			proxyRequest(upstream, c)
			return
		}
	}

	//_ = c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
}

func getUserName(c *gin.Context) string {
	cookieUsername, err := c.Request.Cookie("username")
	if err != nil {
		log.Println("No username was found in cookie + ", err)
	}
	result, err := valkeyServ.Get(c, cookieUsername.Value)
	if err != nil {
		log.Println("No username was found in valkey + ", err)
		redirectToAuth(c)
	}
	return result
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
		req.Header.Set("X-Forwarded-User", "vr.gorbunov")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func renderAuthPage(c *gin.Context) {
	// Отрендерить и вернуть html-страницу с формой для ввода логина и TOTP-кода
}

func validateTotp(c *gin.Context) {
	// Проверить введенный TOTP. Если все ок - проставить куку и отредиректить на url, указанный в параметре redirectUrl.
	// Если нет - отрендерить ту же форму что и в методе выше, но с сообщением об ошибке.
}
