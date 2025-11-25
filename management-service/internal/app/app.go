package app

import (
	"cousework/internal/db"
)

type App struct {
	db *db.Database
}

func NewApp(db *db.Database) *App {
	return &App{db: db}
}

func (app *App) Init() {
	if err := app.db.InitDb(); err != nil {
		panic(err)
	}
}
