package services

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"net/http"
)

func GenerateTOTP(c *gin.Context, email string) string {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: email,
	})
	if err != nil {
		panic(err)
	}

	secret := key.Secret()
	return secret
}

func ValidateTotp(c *gin.Context) {
	// Проверить введенный TOTP. Если все ок - проставить куку и отредиректить на url, указанный в параметре redirectUrl.
	// Если нет - отрендерить ту же форму что и в методе выше, но с сообщением об ошибке.
}

func RenderAuthPage(c *gin.Context) (string, string) {
	// Отрендерить и вернуть html-страницу с формой для ввода логина и TOTP-кода
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
	email := c.PostForm("email")
	password := c.PostForm("password")

	return email, password
}
