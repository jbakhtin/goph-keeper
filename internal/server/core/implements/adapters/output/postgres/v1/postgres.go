package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

type IConfig interface {
	GetDataBaseDriver() string
	GetDataBaseDSN() string
}

type Postgres struct {
	*sql.DB
}

func New(cfg IConfig) (*Postgres, error) {
	db, err := sql.Open(cfg.GetDataBaseDriver(), cfg.GetDataBaseDSN())
	if err != nil {
		return nil, errors.Wrap(err, "db open")
	}

	fmt.Println(cfg.GetDataBaseDriver(), cfg.GetDataBaseDSN())

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

	return &Postgres{
		DB: db,
	}, nil
}
