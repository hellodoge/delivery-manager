package repository

import (
	deliveryManager "github.com/hellodoge/delivery-manager"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"path/filepath"
)

const (
	queriesFolder = "queries"
)

type Authorization interface {
	CreateUser(user deliveryManager.User) (int, error)
	GetUser(username string) (deliveryManager.User, error)
	GetUserByID(id int) (deliveryManager.User, error)
}

type DMProduct interface {
	Create(product deliveryManager.DMProduct) (int, error)
	Search(query deliveryManager.DMProductSearchQuery) ([]deliveryManager.DMProduct, error)
	Exists(productID int) (bool, error)
}

type DMList interface {
	Create(userId int, list deliveryManager.DMList) (int, error)
	GetUserLists(userId int) ([]deliveryManager.DMList, error)
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
		DMProduct: NewDMProductPostgres(db),
		DMList: NewDMListPostgres(db),
	}
}

type DatabaseConfig struct {
	Host string
	Port uint16
	Username string
	Password string
	DBName string
	SSLMode string
}

func getQuery(folder, filename string) (string, error) {
	path := filepath.Join(queriesFolder, folder, filename)
	query, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(query), nil
}