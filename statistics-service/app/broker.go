package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"statistics-service/db"
)

type Consumer struct {
	conn *kafka.Reader
	app  *App
}

func NewConsumer(conn *kafka.Reader, app *App) *Consumer {
	return &Consumer{conn: conn, app: app}
}

func (c *Consumer) Listen() {
	for {
		msg, err := c.conn.ReadMessage(context.TODO())
		if err != nil {
			slog.Error("Не удалось получить сообщение из очереди: ", err.Error())
			return
		}

		err = c.conn.CommitMessages(context.TODO(), msg)
		if err != nil {
			slog.Error("Ошибка коммита сообщения в kafka: ", err.Error())
			continue
		}

		var unmarshalledMessage db.ProcessStatsInput
		err = json.Unmarshal(msg.Value, &unmarshalledMessage)
		if err != nil {
			slog.Error("Не удалось размаршаллить сообщение: ", err.Error())
			continue
		}

		slog.Info(fmt.Sprintf("Получено сообщение: %+v", unmarshalledMessage))

		err = c.app.ProcessStats(unmarshalledMessage)
		if err != nil {
			slog.Error("Ошибка обработки сообщения из Kafka: ", err.Error())
			continue
		}
	}
}
