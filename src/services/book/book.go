package book

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/michaelahli/octopus/src/storage"
	"github.com/michaelahli/octopus/src/storage/postgres"
)

type BookSvc struct {
	db *postgres.Storage
}

type Book interface {
	HandleBooks(w http.ResponseWriter, r *http.Request)
}

func New(db *postgres.Storage) *BookSvc {
	return &BookSvc{db: db}
}

func (s *BookSvc) HandleBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	toInt64 := func(s string) int64 {
		num, _ := strconv.Atoi(s)
		return int64(num)
	}
	switch r.Method {
	case http.MethodGet:
		books, count, err := s.db.GetBooks(ctx, storage.FilterBooks{
			Title: r.URL.Query()["title"],
			Pagination: storage.Pagination{
				Limit:  toInt64(r.URL.Query().Get("limit")),
				Offset: toInt64(r.URL.Query().Get("offset")),
			},
		})
		if err != nil {
			http.Error(w, "books not found", http.StatusNotFound)
			return
		}
		_ = count
		b, err := json.Marshal(books)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(b))
		return
	case http.MethodPost:
		books := make([]storage.Book, 0)
		if err := json.NewDecoder(r.Body).Decode(&books); err != nil {
			http.Error(w, fmt.Sprintf("invalid input: %s", err), http.StatusUnprocessableEntity)
			return
		}
		books, err := s.db.CreateBook(ctx, books)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to create book: %s", err.Error()), http.StatusBadRequest)
			return
		}
		b, err := json.Marshal(books)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(b))
		return
	default:
		fmt.Fprint(w, "Unimplemented")
	}
}
