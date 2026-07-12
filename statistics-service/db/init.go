package db

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
)

var (
	//go:embed migrations/*
	migrations embed.FS
)

func PerformMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	return nil
}

func ConnectDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, fmt.Errorf("ошибка коннекта к БД: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка пинга БД: %w", err)
	}

	return db, nil
}
