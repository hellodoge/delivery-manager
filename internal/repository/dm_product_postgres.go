package repository

import (
	deliveryManager "github.com/hellodoge/delivery-manager"
	"github.com/jmoiron/sqlx"
)

const (
	postgresCreateProduct    = "CreateProduct.sql"
	postgresSearchForProduct = "SearchForProduct.sql"
	postgresProductExists	 = "ProductExists.sql"
)

type DMProductPostgres struct {
	db *sqlx.DB
}

func NewDMProductPostgres(db *sqlx.DB) *DMProductPostgres {
	return &DMProductPostgres{
		db: db,
	}
}

func (r *DMProductPostgres) Create(product deliveryManager.DMProduct) (int, error) {
	query, err := getQuery(postgresQueriesFolder, postgresCreateProduct)
	if err != nil {
		return -1, err
	}
	row := r.db.QueryRow(query, product.Title, product.Description, product.Price)

	var id int
	err = row.Scan(&id)
	return id, err
}

func (r *DMProductPostgres) Search(searchQuery deliveryManager.DMProductSearchQuery) ([]deliveryManager.DMProduct, error) {
	query, err := getQuery(postgresQueriesFolder, postgresSearchForProduct)
	if err != nil {
		return nil, err
	}

	var found []deliveryManager.DMProduct
	err = r.db.Select(&found, query, searchQuery.MatchAllFields, searchQuery.Title,
		searchQuery.Description, searchQuery.TitleOrDescription)

	return found, err
}

func (r *DMProductPostgres) Exists(productID int) (bool, error) {

	queryProductExists, err := getQuery(postgresQueriesFolder, postgresProductExists)
	if err != nil {
		return false, err
	}

	row := r.db.QueryRow(queryProductExists, productID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}