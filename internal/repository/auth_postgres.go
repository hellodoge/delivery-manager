package repository

import (
	"database/sql"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	postgresCreateUser             = "CreateUser.sql"
	postgresGetUser                = "GetUser.sql"
	postgresGetUserByID            = "GetUserByID.sql"
	postgresCreateRefreshToken     = "CreateRefreshToken.sql"
	postgresGetUserByRefreshToken  = "GetUserByRefreshToken.sql"
	postgresGetUserRefreshTokens   = "GetUserRefreshTokens.sql"
	postgresInvalidateRefreshToken = "InvalidateRefreshToken.sql"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user dm.User) (int, error) {
	query, err := getQuery(postgresQueriesFolder, postgresCreateUser)
	if err != nil {
		return -1, err
	}
	row := r.db.QueryRow(query, user.Name, user.Username, user.PasswordHash, user.PasswordSalt)
	var id int
	if err = row.Scan(&id); err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Code.Name() {
			case "unique_violation":
				return -1, ErrUserExists
			}
		}
		return -1, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (dm.User, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetUser)
	if err != nil {
		return dm.User{}, err
	}
	var user dm.User
	err = r.db.Get(&user, query, username)
	return user, err
}

func (r *AuthPostgres) GetUserByID(id int) (dm.User, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetUserByID)
	if err != nil {
		return dm.User{}, err
	}
	var user dm.User
	err = r.db.Get(&user, query, id)
	return user, err
}

func (r *AuthPostgres) CreateRefreshToken(userID int, token string, expiresAt time.Time, ip string) error {
	query, err := getQuery(postgresQueriesFolder, postgresCreateRefreshToken)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, token, userID, ip, expiresAt)
	return err
}

func (r *AuthPostgres) GetUserByRefreshToken(token string) (dm.User, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetUserByRefreshToken)
	if err != nil {
		return dm.User{}, err
	}

	var user dm.User
	err = r.db.Get(&user, query, token)
	if err == sql.ErrNoRows {
		return dm.User{}, ErrRefreshTokenNotFound
	}
	return user, err
}

func (r *AuthPostgres) GetUserRefreshTokens(userID int, issuedAfter time.Time) ([]dm.RefreshTokenInfo, error) {
	query, err := getQuery(postgresQueriesFolder, postgresGetUserRefreshTokens)
	if err != nil {
		return nil, err
	}

	var refreshTokens []dm.RefreshTokenInfo
	err = r.db.Select(&refreshTokens, query, userID, issuedAfter)
	return refreshTokens, err
}

func (r *AuthPostgres) InvalidateRefreshTokens(ids []int, userID int) error {
	query, err := getQuery(postgresQueriesFolder, postgresInvalidateRefreshToken)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, id := range ids {
		_, err := tx.Exec(query, id, userID)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				logrus.Error(err2)
			}
			return err
		}
	}

	return tx.Commit()
}
