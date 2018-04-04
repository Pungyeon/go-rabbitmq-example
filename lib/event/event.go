package event

import (
	"github.com/streadway/amqp"
)

func getQueueName() string {
	return "task_queue"
}

func getExchangeName() string {
	return "logs_topic"
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		getQueueName(), // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		getExchangeName(), // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
}
