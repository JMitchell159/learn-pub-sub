package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("error while marshalling struct into json: %v", err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bytes,
	})
	if err != nil {
		return fmt.Errorf("error while publishing json data: %v", err)
	}

	return nil
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return fmt.Errorf("error while encoding struct into gob: %v", err)
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buff.Bytes(),
	})
	if err != nil {
		return fmt.Errorf("error while publishing gob data: %v", err)
	}

	return nil
}

func PublishGameLog(ch *amqp.Channel, msg, userName string) error {
	return PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+userName, routing.GameLog{
		CurrentTime: time.Now(),
		Message:     msg,
		Username:    userName,
	})
}
