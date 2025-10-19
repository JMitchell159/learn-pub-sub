package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("error while creating a connection: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connection successful!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error while creating a channel from connection: %v", err)
	}

	err = pubsub.SubscribeGob(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable, func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	})
	if err != nil {
		log.Fatalf("error while subscribing to game logs: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "pause":
			log.Println("sending pause message")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Printf("error while publishing pause data: %v", err)
			}
		case "resume":
			log.Println("sending a resume message")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Printf("error while publishing resume data: %v", err)
			}
		case "quit":
			log.Println("exiting server")
			fmt.Println("Shutting down Peril...")
			os.Exit(0)
		default:
			log.Println("unknown command")
		}
	}
}
