package repository

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"
)

func TestDMListPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	var repo = DMListPostgres{db: db}

	err = os.Chdir(path.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}

	const TestTitle = "foo"
	const TestDescription = "bar"
	const TestUserID = 1
	const TestListID = 2
	const TestMeaninglessValue = 3

	tests := []struct {
		name          string
		userID        int
		list          dm.List
		listID        int
		mockBehaviour func(t *testing.T, mock sqlmock.Sqlmock, list dm.List)
		wantError     bool
		wantedError   error
	}{
		{
			name:   "OK",
			userID: TestUserID,
			list: dm.List{
				Title:       TestTitle,
				Description: TestDescription,
			},
			listID: TestListID,
			mockBehaviour: func(t *testing.T, mock sqlmock.Sqlmock, list dm.List) {
				mock.ExpectBegin()
				query, err := ioutil.ReadFile(path.Join(queriesFolder,
					postgresQueriesFolder, postgresCreateList))
				if err != nil {
					t.Error(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).AddRow(TestListID)
				mock.ExpectQuery(regexp.QuoteMeta(string(query))).WithArgs(TestTitle,
					TestDescription).WillReturnRows(rows)

				query, err = ioutil.ReadFile(path.Join(queriesFolder,
					postgresQueriesFolder, postgresLinkListWithUser))
				if err != nil {
					t.Error(err)
				}
				mock.ExpectExec(regexp.QuoteMeta(string(query))).WithArgs(TestUserID,
					TestListID).WillReturnResult(sqlmock.NewResult(TestMeaninglessValue, TestMeaninglessValue))

				mock.ExpectCommit()
			},
		},
		{
			name:   "Database failed (while creating list)",
			userID: TestUserID,
			list: dm.List{
				Title:       TestTitle,
				Description: TestDescription,
			},
			mockBehaviour: func(t *testing.T, mock sqlmock.Sqlmock, list dm.List) {
				mock.ExpectBegin()
				query, err := ioutil.ReadFile(path.Join(queriesFolder,
					postgresQueriesFolder, postgresCreateList))
				if err != nil {
					t.Error(err)
				}
				mock.ExpectQuery(regexp.QuoteMeta(string(query))).WithArgs(TestTitle,
					TestDescription).WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name:   "Database failed (while linking list with user)",
			userID: TestUserID,
			list: dm.List{
				Title:       TestTitle,
				Description: TestDescription,
			},
			mockBehaviour: func(t *testing.T, mock sqlmock.Sqlmock, list dm.List) {
				mock.ExpectBegin()
				query, err := ioutil.ReadFile(path.Join(queriesFolder,
					postgresQueriesFolder, postgresCreateList))
				if err != nil {
					t.Error(err)
				}
				rows := sqlmock.NewRows([]string{"id"}).AddRow(TestListID)
				mock.ExpectQuery(regexp.QuoteMeta(string(query))).WithArgs(TestTitle,
					TestDescription).WillReturnRows(rows)

				query, err = ioutil.ReadFile(path.Join(queriesFolder,
					postgresQueriesFolder, postgresLinkListWithUser))
				if err != nil {
					t.Error(err)
				}
				mock.ExpectExec(regexp.QuoteMeta(string(query))).WithArgs(TestUserID,
					TestListID).WillReturnError(errors.New("db error"))

				mock.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehaviour(t, mock, test.list)
		})

		result, err := repo.Create(test.userID, test.list)
		if test.wantError {
			if test.wantedError != nil {
				assert.Equal(t, err, test.wantedError)
			} else {
				if err == nil {
					t.Error("wanted error, got nil")
				}
			}
		} else if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, result, test.listID)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	}
}
