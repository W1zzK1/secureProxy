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
	secret := "PDSGYV7LJSFSPPRTSJZCBHVV3RHXCKZV"

	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return
	}
	println("Code:", code)

	validationResult := totp.Validate(code, secret)
	println("Validation passed :", validationResult)
}

/* ERROR
# github.com/valkey-io/valkey-glide/go/v2/appConfig
..\..\..\go\pkg\mod\github.com\valkey-io\valkey-glide\go\v2@v2.1.0\appConfig\pubsub_subscription_config.go:11:43: undefined: models.PubSubMessage
*/
//func Test_valkeySlide(t *testing.T) {
//	host := "localhost"
//	port := 6379
//
//	appConfig := appConfig.NewClientConfiguration().
//		WithAddress(&appConfig.NodeAddress{Host: host, Port: port})
//
//	client, err := glide.NewClient(appConfig)
//	if err != nil {
//		panic(err)
//	}
//	defer client.Close()
//
//	context := context.Background()
//	res, err := client.Ping(context)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(res) // PONG
//}

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
