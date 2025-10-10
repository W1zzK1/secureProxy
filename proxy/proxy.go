package proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"secureProxy/services"
)

func HandleAuthDomain(c *gin.Context, valkeyServ *services.ValkeyService) {
	if c.Request.Method == "GET" {
		// Получаем redirectUrl из query параметра
		redirectUrl := c.Query("redirectUrl")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"RedirectUrl": redirectUrl, // Передаем в шаблон
		})
		return
	}

	username := c.PostForm("Username")
	totp := c.PostForm("TOTP")
	redirectUrl := c.PostForm("redirectUrl") // Получаем из скрытого поля формы

	if username == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error":       "username is required",
			"RedirectUrl": redirectUrl, // Сохраняем при ошибке
		})
		return
	}
	if totp == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error":       "Totp is required",
			"RedirectUrl": redirectUrl, // Сохраняем при ошибке
		})
		return
	}

	validation := services.ValidateTotp(c, totp, "OWX2WB6TEUBMPYMSXML4B2YKFEEQ5FYI")

	if !validation {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"Error":       "Invalid TOTP",
			"RedirectUrl": redirectUrl, // Сохраняем при ошибке
		})
		return
	}

	// Создаем сессию
	sessionCookie := "session:" + SetProxyCookie(c)
	valkeyServ.Set(c, sessionCookie, username)
	valkeyServ.Expire(c, sessionCookie, 1800)

	// Используем правильный статус код для редиректа - 302
	if redirectUrl != "" {
		c.Redirect(http.StatusFound, redirectUrl) // http.StatusFound = 302
	} else {
		// Fallback если redirectUrl пустой
		c.Redirect(http.StatusFound, "https://site1.secure-proxy.lan:8443/")
	}
}

func SetProxyCookie(c *gin.Context) string {
	sessionId := uuid.New().String()
	c.SetCookie("SECURE_PROXY_SESSION", sessionId, 60*60, "/", ".secure-proxy.lan", true, true)
	return sessionId
}
