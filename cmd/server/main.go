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
		log.Printf("error while creating a connection: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connection successful!")

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("error while creating a channel from connection: %v", err)
	}

	gamelogic.PrintServerHelp()

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", pubsub.Durable)
	if err != nil {
		log.Printf("error while declaring and binding game logs: %v", err)
	}

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
