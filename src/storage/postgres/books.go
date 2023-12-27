package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/michaelahli/octopus/src/storage"
	dbutil "github.com/michaelahli/octopus/svcutils/storage/postgres"
	timeutil "github.com/michaelahli/octopus/svcutils/time"
)

const (
	getBooksQuery   = `SELECT * FROM books WHERE %s`
	countBooksQuery = `SELECT COUNT(*) FROM books WHERE %s`
	createBookQuery = `INSERT INTO books (
		title
	) VALUES (
		:title
	) RETURNING *`
)

func (s *Storage) GetBooks(ctx context.Context, req storage.FilterBooks) ([]storage.Book, int64, error) {
	var (
		sort, pagination string
		count            int64
		filters          = make(map[string][]string)
		spQueries        = []dbutil.SpecialQuery{}
		lbs              = make([]storage.Book, 0)
	)

	for key, value := range map[string][]string{
		"category": req.Title,
		"deleted":  {"IS NULL"},
	} {
		filters[key] = value
	}

	if !req.Pagination.From.IsZero() {
		spQueries = append(spQueries, dbutil.SpecialQuery{
			Query: "created > ?",
			Args:  []string{req.Pagination.From.Format(timeutil.FormatYYYYMMDD)},
			Type:  dbutil.And,
		})
	}

	if !req.Pagination.To.IsZero() {
		spQueries = append(spQueries, dbutil.SpecialQuery{
			Query: "created < ?",
			Args:  []string{req.Pagination.To.Format(timeutil.FormatYYYYMMDD)},
			Type:  dbutil.And,
		})
	}

	if len(req.Pagination.Order) > 0 {
		sort = `ORDER BY`
	}

	for _, order := range req.Pagination.Order {
		sort = strings.Join([]string{sort, order.Column, order.Direction}, " ")
	}

	if req.Pagination.Limit > 0 {
		pagination = fmt.Sprintf(`LIMIT %d OFFSET %d`, req.Pagination.Limit, req.Pagination.Offset)
	}

	query, args, err := dbutil.BuildFilter(getBooksQuery, filters, spQueries...)
	if err != nil {
		return nil, count, fmt.Errorf("internal error on building filter: %v", err)
	}

	query = strings.Join([]string{query, sort, pagination}, " ")

	stmt, err := s.prepare(ctx, query)
	if err != nil {
		return nil, count, err
	}
	defer stmt.Close()

	if err = stmt.SelectContext(ctx, &lbs, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, count, storage.ErrNotFound
		}
		return nil, count, err
	}

	query, args, err = dbutil.BuildFilter(countBooksQuery, filters, spQueries...)
	if err != nil {
		return nil, count, fmt.Errorf("internal error on building filter: %v", err)
	}

	cstmt, err := s.prepare(ctx, query)
	if err != nil {
		return nil, count, err
	}
	defer cstmt.Close()

	if err = cstmt.GetContext(ctx, &count, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, count, storage.ErrNotFound
		}
		return nil, count, err
	}

	return lbs, count, nil
}

func (s *Storage) CreateBook(ctx context.Context, reqs []storage.Book) ([]storage.Book, error) {
	ltr := make([]storage.Book, len(reqs))

	stmt, err := s.prepareNamed(ctx, createBookQuery)
	if err != nil {
		return nil, err
	}

	for i, req := range reqs {
		tr := new(storage.Book)
		if err = stmt.GetContext(ctx, tr, req); err != nil {
			return nil, err
		}
		ltr[i] = *tr
	}

	return ltr, nil
}
