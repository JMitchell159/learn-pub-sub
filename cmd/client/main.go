package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		log.Fatalf("error while creating a connection: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error while creating a channel from connection: %v", err)
	}

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error while running client welcome: %v", err)
	}

	state := gamelogic.NewGameState(userName)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", userName), routing.PauseKey, pubsub.Transient, handlerPause(state))
	if err != nil {
		log.Fatalf("error while subscribing to pause queue: %v", err)
	}

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, fmt.Sprintf("army_moves.%s", userName), routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerArmyMoves(state, ch))
	if err != nil {
		log.Fatalf("error while subscribing to army moves queue: %v", err)
	}

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, "war", routing.WarRecognitionsPrefix+"."+userName, pubsub.Durable, handlerWar(state, ch))
	if err != nil {
		log.Fatalf("error while subscribing to war declarations queue: %v", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch input[0] {
		case "spawn":
			err = state.CommandSpawn(input)
		case "move":
			move, err := state.CommandMove(input)
			if err != nil {
				fmt.Println("Move unsuccessful")
			} else {
				fmt.Println("Move successful")
			}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, fmt.Sprintf("army_moves.%s", userName), move)
			if err != nil {
				log.Println("publish move unsuccessful")
			} else {
				log.Println("publish move successful")
			}
		case "status":
			state.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			if len(input) < 2 {
				fmt.Println("You must provide a numerical argument for spam")
				continue
			}
			num, err := strconv.Atoi(input[1])
			if err != nil {
				fmt.Printf("You must provide a numerical argument for spam: %v\n", err)
				continue
			}
			for i := 0; i < num; i++ {
				err := pubsub.PublishGameLog(ch, gamelogic.GetMaliciousLog(), userName)
				if err != nil {
					log.Println("spam publish failed")
				}
			}
		case "quit":
			gamelogic.PrintQuit()
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
		}

		if err != nil {
			log.Println(err)
		}
	}
}
