package main

import (
	_ "cousework_auth/docs"
	"cousework_auth/internal/app"
	"cousework_auth/internal/db"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"log/slog"
	"os"
)

// @title Auth-Service
// @version 2.0
// @description Сервис аутентификации
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
func main() {
	database, err := initDb()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer database.Close()

	slog.Warn("initDb")

	application := app.NewApp(db.NewDatabase(database))
	application.Init()

	slog.Warn("Init")

	server := application.CreateServer(":8081")
	slog.Info("Server started on :8081")

	slog.Warn("CreateServer")

	slog.Info("Auth-service server starting...")
	log.Fatal(server.ListenAndServe())
}

func initDb() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, fmt.Errorf("Ошибка коннекта к БД: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Ошибка пинга БД: %w", err)
	}

	return db, nil
}
