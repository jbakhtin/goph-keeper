package postgres

import (
	"database/sql"
	"embed"
	"github.com/pkg/errors"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config interface {
	GetDataBaseDriver() string
	GetDataBaseDSN() string
}

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

func NewSQLClient(cfg Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.GetDataBaseDriver(), cfg.GetDataBaseDSN())
	if err != nil {
		return nil, errors.Wrap(err, "db open")
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "db ping")
	}

	goose.SetBaseFS(EmbedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, errors.Wrap(err, "set dialect")
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return nil, errors.Wrap(err, "run migrations")
	}

	return db, nil
}
