package main

import (
	"fmt"
	"os"

	"github.com/Pungyeon/go-rabbitmq-example/lib/event"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	emitter, err := event.NewEventEmitter(conn)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 10; i++ {
		emitter.Push(fmt.Sprintf("[%d] - %s", i, os.Args[1]), os.Args[1])
	}
}
