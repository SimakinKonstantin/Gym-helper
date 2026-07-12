package main

import (
	_ "cousework/docs"
	"cousework/internal/app"
	"cousework/internal/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
)

func main() {
	database, err := initDb()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer database.Close()

	application := app.NewApp(db.NewDatabase(database))
	application.Init()

	server := application.CreateServer(":8082")
	slog.Info("Server started on :8082")

	log.Fatal(server.ListenAndServe())
}

func initDb() (*sqlx.DB, error) {
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
