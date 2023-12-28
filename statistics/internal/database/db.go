package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/patrick246/solar/statistics/internal/config"
)

//go:embed migrations
var migrations embed.FS

func Connect(cfg config.Database) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	driver, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		return err
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return err
	}

	version, dirty, err := migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		version = 0
	} else if err != nil {
		return err
	}

	logger.InfoContext(ctx, "before migration", "version", version, "dirty", dirty)

	err = migrator.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.InfoContext(ctx, "no change")
	} else if err != nil {
		return err
	}

	version, dirty, err = migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		version = 0
	} else if err != nil {
		return err
	}

	logger.InfoContext(ctx, "after migration", "version", version, "dirty", dirty)

	return nil
}
