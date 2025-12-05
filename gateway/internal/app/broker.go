package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func (app *App) Write(value ProcessTrainingInput) error {
	marshalled, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Ошибка отправки в kafka")
	}

	return app.producer.WriteMessages(context.TODO(), kafka.Message{
		Value: marshalled,
	})
}
