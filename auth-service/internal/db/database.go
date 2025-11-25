package db

import (
	"cousework_auth/internal"
	"embed"
	_ "embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var (
	//go:embed scripts/add_user.sql
	addUser string

	//go:embed scripts/get_user.sql
	getUser string
)

var (
	//go:embed migrations/*
	embedMigrations embed.FS
)

type Database struct {
	Db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{Db: db}
}

func (d *Database) InitDb() {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(d.Db.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (d *Database) AddUser(login, hash string) error {
	if _, err := d.Db.Exec(addUser, login, hash); err != nil {
		return fmt.Errorf("ошибка сохранения пользователя в БД: %w", err)
	}

	return nil
}

func (d *Database) GetUser(login string) (internal.UserDb, error) {
	var user internal.UserDb
	if err := d.Db.Get(&user, getUser, login); err != nil {
		return internal.UserDb{}, fmt.Errorf("ошибка получения пользователя из БД: %w", err)
	}

	return user, nil
}
