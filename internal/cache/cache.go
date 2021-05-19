package cache

type RefreshTokens interface {
	NewRefreshToken(user *RefreshTokenSavedFields) (string, error)
	GetSavedFields(token string) (*RefreshTokenSavedFields, error)
}

type Storage struct {
	RefreshTokens
}

type StorageConfig struct {
	RTConfig RefreshTokensConfig
}


func NewStorage(options RedisConnectionOptions, config StorageConfig) *Storage {
	return &Storage{
		RefreshTokens: NewRefreshTokensRedis(options, config.RTConfig),
	}
}