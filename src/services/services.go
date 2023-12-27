package services

import (
	"github.com/michaelahli/octopus/src/services/book"
	"github.com/michaelahli/octopus/src/services/common"
	"github.com/michaelahli/octopus/src/storage/postgres"

	"github.com/spf13/viper"
)

type Svc struct {
	*common.CommonSvc
	*book.BookSvc
}

type Services interface {
	common.Common
	book.Book
}

func New(config *viper.Viper, db *postgres.Storage) Services {
	return &Svc{common.New(config), book.New(db)}
}
