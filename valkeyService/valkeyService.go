package valkeyService

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/valkey-io/valkey-go"
	"log"
	"net/http"
)

type ValkeyService struct {
	client valkey.Client
}

func NewValkeyService(client valkey.Client) *ValkeyService {
	return &ValkeyService{client: client}
}

func CreateClient() (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})

	// Проверка подключения
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}
	//defer client.Close()
	return client, nil
}

// Общий паттерн:
//result := client.Do(ctx, client.B().Команда().Параметры().Build()).МетодРазбора()

// Set сохраняет значение по ключу
func (s *ValkeyService) Set(ctx context.Context, key, value string) error {
	err := s.client.Do(ctx, s.client.B().Set().Key(key).Value(value).Build()).Error()
	if err != nil {
		return fmt.Errorf("set failed for key %s: %w", key, err)
	}
	return nil
}

func (s *ValkeyService) SetValue(c *gin.Context) {
	// /valkey/set-value?key=test&value=hello%20world
	ctx := c.Request.Context()

	key := c.Query("key")
	value := c.Query("value")

	if key == "" || value == "" {
		c.JSON(400, gin.H{
			"error": "key or value is empty",
		})
		return
	}

	err := s.Set(ctx, key, value)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to set value for key " + key + " and value " + value + ": " + err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"key":   key,
			"value": value,
		},
	})
}

func (s *ValkeyService) Get(ctx context.Context, key string) (string, error) {
	result, err := s.client.Do(ctx, s.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		log.Fatal("get failed")
	}
	return result, nil
}

func (s *ValkeyService) GetValue(c *gin.Context) {
	// /valkey/get-value/get/test
	ctx := c.Request.Context()
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "key parameter is required",
		})
		return
	}
	value, err := s.Get(ctx, key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "failed to get value for key " + key + ": " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"key":   key,
			"value": value,
		},
	})
}
