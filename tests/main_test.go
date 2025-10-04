package tests

import (
	"context"
	"github.com/pquerna/otp/totp"
	"secureProxy/valkeyService"
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
	secret := "4I7VBTB2EH4DMGKINHTFAY2VDTKO543V"

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return
	}
	println("Code:", code)

	validationResult := totp.Validate(code, secret)
	println("Validation passed :", validationResult)
}

func Test_valkey(t *testing.T) {
	client, err := valkeyService.CreateClient()
	if err != nil {
		panic(err)
	}
	valkeyService := valkeyService.NewValkeyService(client)
	defer client.Close()

	ctx := context.Background()
	valkeyService.Set(ctx, "nickname", "w1zzk1")
	value, err := valkeyService.Get(ctx, "nickname")
	if err != nil {
		return
	}
	println(value)

}
