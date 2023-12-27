package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/michaelahli/octopus/svcutils/storage/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/spf13/viper"

	_ "github.com/lib/pq" // register postgres driver
)

type Storage struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func New(config *viper.Viper) (*Storage, error) {
	db, err := postgres.Connectx(config)
	if err != nil {
		log.Printf("db connection: %s\n", err.Error())
		return nil, err
	}

	log.Printf("Successfully established connection to database\n")
	log.Printf("Ping postgres client : %v \n", db.PingContext(context.Background()))

	return NewDB(db), nil
}

func NewStorageDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func NewTestStorage(dbstring, migrationDir string) (*Storage, func()) {
	db, teardown := postgres.MustNewDevelopmentDB(dbstring, migrationDir)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)
	return NewStorageDB(db), teardown
}

func (s *Storage) RunMigration(dir string) error {
	return goose.Run("up", s.db.DB, dir)
}

type pgTx struct{}

type tx struct {
	*sqlx.Tx
	committed *bool
}

func (s *Storage) NewTransacton(ctx context.Context) (context.Context, error) {
	t, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, pgTx{}, &tx{
		Tx:        t,
		committed: new(bool),
	}), nil
}

func (s *Storage) Commit(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	if err := t.Commit(); err != nil {
		return err
	}
	*t.committed = true
	return nil
}

func (s *Storage) Rollback(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	return t.Rollback()
}

func getTx(ctx context.Context) *tx {
	if t, ok := ctx.Value(pgTx{}).(*tx); ok && !*t.committed {
		return t
	}
	return nil
}

// PrepareNamed prepares a named query in the current transaction (if begun) or with the db.
func (s *Storage) prepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.PrepareNamedContext(ctx, query)
	}
	return s.db.PrepareNamedContext(ctx, query)
}

// Prepare prepares a query in the current transaction (if begun) or with the db.
func (s *Storage) prepare(ctx context.Context, query string) (*sqlx.Stmt, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.PreparexContext(ctx, query)
	}
	return s.db.PreparexContext(ctx, query)
}
