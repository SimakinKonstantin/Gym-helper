package app

import (
	"cousework_auth/internal/db"
	_ "embed"
)

type App struct {
	db *db.Database
}

func NewApp(db *db.Database) *App {
	return &App{db: db}
}

func (app *App) Init() {
	app.db.InitDb()
}
