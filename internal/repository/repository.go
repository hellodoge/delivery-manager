package repository

import (
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"path/filepath"
	"time"
)

const (
	queriesFolder = "queries"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user dm.User) (int, error)
	GetUser(username string) (dm.User, error)
	GetUserByID(id int) (dm.User, error)
	CreateRefreshToken(userID int, token string, expiresAt time.Time, ip string) error
	GetUserByRefreshToken(token string) (dm.User, error)
	GetUserRefreshTokens(userID int, issuedAfter time.Time) ([]dm.RefreshTokenInfo, error)
	InvalidateRefreshTokens(ids []int, userID int) error
	GetActiveRefreshTokens(userID int) ([]dm.RefreshTokenInfo, error)
}

type DMProduct interface {
	Create(product []dm.Product) ([]int, error)
	Search(query dm.ProductSearchQuery) ([]dm.Product, error)
	Exists(productID int) (bool, error)
}

type DMList interface {
	Create(userId int, list dm.List) (int, error)
	GetUserLists(userId int) ([]dm.List, error)
	GetOwners(listId int) ([]int, error)
	Delete(listId int) error
	AddProduct(listID, productID, count int) error
}

type DMDelivery interface {
}

type Repository struct {
	Authorization
	DMProduct
	DMList
	DMDelivery
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		DMProduct:     NewDMProductPostgres(db),
		DMList:        NewDMListPostgres(db),
	}
}

type DatabaseConfig struct {
	Host     string
	Port     uint16
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func getQuery(folder, filename string) (string, error) {
	path := filepath.Join(queriesFolder, folder, filename)
	query, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(query), nil
}
