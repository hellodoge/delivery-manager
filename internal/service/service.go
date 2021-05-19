package service

import (
	deliveryManager "github.com/hellodoge/delivery-manager"
	"github.com/hellodoge/delivery-manager/internal/cache"
	"github.com/hellodoge/delivery-manager/internal/repository"
	"time"
)

const (
	DefaultTokenLifetime = 12 * time.Hour
)

type Authorization interface {
	CreateUser(user deliveryManager.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type DMProduct interface {
	Create(products []deliveryManager.DMProduct) ([]deliveryManager.DMProduct, error)
	Search(query deliveryManager.DMProductSearchQuery) ([]deliveryManager.DMProduct, error)
	Exists(productID int) (bool, error)
}

type DMList interface {
	Create(userID int, list deliveryManager.DMList) (deliveryManager.DMList, error)
	GetUserLists(userID int) ([]deliveryManager.DMList, error)
	Delete(userID, listID int) error
	AddProduct(userID, listID int, index []deliveryManager.DMProductIndex) error
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
	TokenLifetime time.Duration
	CheckHash     bool
}

func NewService(repo *repository.Repository, cache *cache.Storage, config Config) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, cache.RefreshTokens, config.AuthConfig),
		DMProduct: NewDMProductService(repo.DMProduct),
		DMList: NewDMListService(repo.DMList),
	}
}