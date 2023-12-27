package services

import (
	"github.com/michaelahli/octopus/src/services/book"
	"github.com/michaelahli/octopus/src/services/common"

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

func New(config *viper.Viper) Services {
	return &Svc{common.New(config), book.New(config)}
}
