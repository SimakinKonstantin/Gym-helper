package main

import (
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"log"
	"log/slog"
	"os"
	"statistics-service/app"
	"statistics-service/db"
	_ "statistics-service/docs"
)

func main() {
	connection, err := db.ConnectDb()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer connection.Close()

	err = db.PerformMigrations(connection)
	if err != nil {
		log.Fatal(err)
		return
	}

	database := db.NewDatabase(connection)
	application := app.NewApp(database)

	consumerConn := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{os.Getenv("KAFKA_ADDR")}, Topic: os.Getenv("KAFKA_TOPIC"), GroupID: "group-id", CommitInterval: 0})
	defer consumerConn.Close()

	consumer := app.NewConsumer(consumerConn, application)

	go consumer.Listen()

	server := application.CreateServer(":8083")
	slog.Info("Server started on :8083")

	log.Fatal(server.ListenAndServe())
}
