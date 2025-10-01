package tests

import (
	"github.com/pquerna/otp/totp"
	"testing"
	"time"
)

func Test_generateTOTP(t *testing.T) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
	})
	if err != nil {
		panic(err)
	}

	secret := key.Secret()
	println(secret)
}

func Test_validateTOTP(t *testing.T) {
	secret := "4NKVI67WZ7HZOWKJFLMUASG4DFU5O3IO"

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return
	}
	println("Code:", code)

	validationResult := totp.Validate(code, secret)
	println("Validation passed :", validationResult)
}
