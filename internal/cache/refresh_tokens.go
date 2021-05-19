package cache

import "time"

type RefreshTokensConfig struct {
	Expiration time.Duration
	DB int
}

type RefreshTokenSavedFields struct {
	UserID int `json:"user_id"`
	UserHashPart string `json:"user_hash_part"`
}