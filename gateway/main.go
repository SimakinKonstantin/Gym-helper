package main

import (
	"coursework_gateway/internal/app"
	"github.com/segmentio/kafka-go"
	"log"
	"log/slog"
	"os"
)

func main() {
	producer := kafka.Writer{Addr: kafka.TCP(os.Getenv("KAFKA_ADDR")), Topic: os.Getenv("KAFKA_TOPIC")}
	defer producer.Close()

	application := app.NewApp(&producer)

	server := application.CreateServer(":8080")
	slog.Info("Server started on :8080")
	log.Fatal(server.ListenAndServe())
}
