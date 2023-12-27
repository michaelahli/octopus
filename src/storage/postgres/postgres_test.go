package postgres

import (
	"log"

	"github.com/jmoiron/sqlx"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func NewMock() (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
