package book

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type BookSvc struct {
	config *viper.Viper
}

type Book interface {
	HandleBooks(w http.ResponseWriter, r *http.Request)
}

func New(config *viper.Viper) *BookSvc {
	return &BookSvc{config: config}
}

func (s *BookSvc) HandleBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Unimplemented")
}
