package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error while creating channel: %v", err)
	}

	queue, err := ch.QueueDeclare(queueName, queueType == Durable, queueType == Transient, queueType == Transient, false, nil)
	if err != nil {
		return ch, amqp.Queue{}, fmt.Errorf("error while declaring queue: %v", err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return ch, queue, fmt.Errorf("error while binding queue to channel: %v", err)
	}

	return ch, queue, nil
}
