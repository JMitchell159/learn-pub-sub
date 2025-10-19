package pubsub

import (
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

func Subscribe[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType, decoder func([]byte) (T, error)) error {
	var v T
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("error while declaring and binding queue: %v", err)
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("error while setting prefetch: %v", err)
	}

	delivery, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error while consuming channel: %v", err)
	}

	go func() {
		for d := range delivery {
			v, err = decoder(d.Body)
			if err != nil {
				log.Printf("error while unmarshaling delivery into struct: %v", err)
			}

			ack := handler(v)

			switch ack {
			case Ack:
				err = d.Ack(false)
				if err != nil {
					log.Printf("error while acknowledging delivery: %v", err)
				}
			case NackRequeue:
				err = d.Nack(false, true)
				if err != nil {
					log.Printf("error while negative acknowledging and requeueing: %v", err)
				}
			case NackDiscard:
				err = d.Nack(false, false)
				if err != nil {
					log.Printf("error while negative acknowledging and discarding: %v", err)
				}
			}
		}
	}()

	return nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	return Subscribe(conn, exchange, queueName, key, queueType, handler, decodeJSON)
}

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
	return Subscribe(conn, exchange, queueName, key, queueType, handler, decodeGob)
}
