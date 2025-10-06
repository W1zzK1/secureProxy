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

func ValidateTotp(c *gin.Context, code, secret string) bool {
	return totp.Validate(code, secret)
}

func RenderTOTPPage(c *gin.Context, secret string) bool {
	// Если это GET запрос - рендерим форму
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "totpValidate.html", gin.H{
			"Secret": secret,
		})
		return false
	}

	// Если это POST запрос - получаем код
	code := c.PostForm("code")
	return totp.Validate(code, secret)
}
