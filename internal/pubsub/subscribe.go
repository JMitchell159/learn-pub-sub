package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	var v T
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("error while declaring and binding queue: %v", err)
	}

	delivery, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error while consuming channel: %v", err)
	}

	go func() {
		for d := range delivery {
			err = json.Unmarshal(d.Body, &v)
			if err != nil {
				log.Printf("error while unmarshaling delivery into struct: %v", err)
			}

			ack := handler(v)

			switch ack {
			case Ack:
				log.Println("Ack")
				err = d.Ack(false)
				if err != nil {
					log.Printf("error while acknowledging delivery: %v", err)
				}
			case NackRequeue:
				log.Println("Nack Requeue")
				err = d.Nack(false, true)
				if err != nil {
					log.Printf("error while negative acknowledging and requeueing: %v", err)
				}
			case NackDiscard:
				log.Println("Nack Discard")
				err = d.Nack(false, false)
				if err != nil {
					log.Printf("error while negative acknowledging and discarding: %v", err)
				}
			}
		}
	}()

	return nil
}
