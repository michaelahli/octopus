package storage

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type Pagination struct {
	Limit  int64
	Offset int64
	Order  []Sort
	From   time.Time
	To     time.Time
}

type Sort struct {
	Column    string
	Direction string
}

type FilterBooks struct {
	Title      []string
	Pagination Pagination
}

type Book struct {
	BookId  string     `db:"book_id" json:"book_id"`
	Title   string     `db:"title" json:"title"`
	Created time.Time  `db:"created" json:"created"`
	Updated *time.Time `db:"updated" json:"updated"`
	Deleted *time.Time `db:"deleted" json:"deleted"`
}
