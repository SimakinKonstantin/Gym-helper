package app

import "github.com/segmentio/kafka-go"

type App struct {
	producer *kafka.Writer
}

func NewApp(producer *kafka.Writer) *App {
	return &App{producer}
}
