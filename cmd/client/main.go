package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Printf("error while creating a connection: %v", err)
	}
	defer conn.Close()

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("error while running client welcome: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", userName), routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Printf("error while running declare and bind: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down Peril Client...")
}
