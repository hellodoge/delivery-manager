package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/hellodoge/delivery-manager/pkg/auth"
	"github.com/sirupsen/logrus"
)

const (
	TimesToTry=10
)

type RefreshTokensRedis struct {
	client *redis.Client
	config RefreshTokensConfig
}

func NewRefreshTokensRedis(connectionOptions RedisConnectionOptions, config RefreshTokensConfig) *RefreshTokensRedis {
	client, err := RedisConnect(connectionOptions, config.DB)
	if err != nil {
		logrus.Fatal(err)
	}
	return &RefreshTokensRedis{
		client: client,
		config: config,
	}
}

func (c *RefreshTokensRedis) NewRefreshToken(user *RefreshTokenSavedFields) (string, error) {
	var token string
	serialized, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	for i := 0; i < TimesToTry; i++ {
		token, err = auth.GenerateRefreshToken(c.config.TokenBytesLength)
		if err != nil {
			return "", err
		}
		set, err := c.client.SetNX(context.Background(), token, serialized, c.config.Expiration).Result()
		if err != nil {
			return "", err
		}
		if set {
			return token, nil
		} else {
			logrus.Warn("Created refresh token collided with existing (%s)", token)
		}
	}
	return "", errors.New("max number of tries reached while creating a token")
}

func (c *RefreshTokensRedis) GetSavedFields(token string) (*RefreshTokenSavedFields, error) {
	data, err := c.client.Get(context.Background(), token).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var savedFields = new(RefreshTokenSavedFields)
	err = json.Unmarshal([]byte(data), &savedFields)
	if err != nil {
		return nil, err
	}

	return savedFields, nil
}