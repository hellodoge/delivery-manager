package service

import (
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/cache"
	"github.com/hellodoge/delivery-manager/internal/repository"
	"time"
)

const (
	DefaultTokenLifetime = 12 * time.Hour
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user dm.User) (int, error)
	GenerateToken(refreshToken string) (string, error)
	ParseToken(token string) (int, error)
	GenerateRefreshToken(username, password string, ip string) (string, error)
	GetUserRefreshTokens(userID int, issuedAfterString string) ([]dm.RefreshTokenInfo, error)
	InvalidateRefreshTokens(ids []int, userID int) error
	GetActiveRefreshTokens(userID int) ([]dm.RefreshTokenInfo, error)
}

type DMProduct interface {
	Create(products []dm.Product) ([]dm.Product, error)
	Search(query dm.ProductSearchQuery) ([]dm.Product, error)
	Exists(productID int) (bool, error)
}

type DMList interface {
	Create(userID int, list dm.List) (dm.List, error)
	GetUserLists(userID int) ([]dm.List, error)
	Delete(userID, listID int) error
	AddProduct(userID, listID int, index []dm.ProductIndex) error
}

type DMDelivery interface {
}

type Service struct {
	Authorization
	DMProduct
	DMList
	DMDelivery
}

type Config struct {
	AuthConfig AuthServiceConfig
}

type AuthServiceConfig struct {
	TokenLifetime        time.Duration
	RefreshTokenLifetime time.Duration
	CheckHash            bool
}

func NewService(repo *repository.Repository, cache *cache.Storage, config Config) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, cache.RefreshTokens, config.AuthConfig),
		DMProduct:     NewDMProductService(repo.DMProduct),
		DMList:        NewDMListService(repo.DMList),
	}
}
