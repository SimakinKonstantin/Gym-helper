package app

import "statistics-service/db"

type App struct {
	db *db.Database
}

func NewApp(db *db.Database) *App {
	return &App{db}
}
