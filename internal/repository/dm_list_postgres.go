package repository

import (
	"fmt"
	deliveryManager "github.com/hellodoge/delivery-manager"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	postgresCreateList       = "CreateList.sql"
	postgresLinkListWithUser = "LinkListWithUser.sql"
	postgresGetUserLists     = "GetUserLists.sql"
	postgresGetListOwners    = "GetListOwners.sql"
	postgresDeleteList       = "DeleteList.sql"
	postgresAddProduct		 = "AddProduct.sql"
)

type DMListPostgres struct {
	db *sqlx.DB
}

func NewDMListPostgres(db *sqlx.DB) *DMListPostgres {
	return &DMListPostgres{
		db: db,
	}
}

func (r *DMListPostgres) Create(userId int, list deliveryManager.DMList) (int, error) {

	queryCreateList, err := getQuery(postgresQueriesFolder, postgresCreateList)
	if err != nil {
		return -1, err
	}

	queryLinkListWithUser, err := getQuery(postgresQueriesFolder, postgresLinkListWithUser)
	if err != nil {
		return -1, err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	row := tx.QueryRow(queryCreateList, list.Title, list.Description)
	var listId int
	if err := row.Scan(&listId); err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			logrus.Error(err2)
		}
		return -1, fmt.Errorf("cannot create list: %s", err)
	}

	_, err = tx.Exec(queryLinkListWithUser, userId, listId)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			logrus.Error(err2)
		}
		return -1, fmt.Errorf("cannot link list with user: %s", err)
	}

	return listId, tx.Commit()
}

func (r *DMListPostgres) GetUserLists(userId int) ([]deliveryManager.DMList, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetUserLists)
	if err != nil {
		return nil, err
	}

	var lists []deliveryManager.DMList
	err = r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *DMListPostgres) GetOwners(listId int) ([]int, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetListOwners)
	if err != nil {
		return nil, err
	}

	var owners []int
	err = r.db.Select(&owners, query, listId)

	return owners, err
}

func (r *DMListPostgres) Delete(listId int) error {
	query, err := getQuery(postgresQueriesFolder, postgresDeleteList)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(query, listId)
	return err
}


func (r *DMListPostgres) AddProduct(listID, productID, count int) error {
	queryAddProduct, err := getQuery(postgresQueriesFolder, postgresAddProduct)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(queryAddProduct, listID, productID, count)

	return err
}