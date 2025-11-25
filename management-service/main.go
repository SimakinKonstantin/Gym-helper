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

// @title Auth-Service
// @version 2.0
// @description Сервис аутентификации
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
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
