package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateConnection(url string) *amqp.Channel {
	connection, err := amqp.Dial(url)

	if err != nil {
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	return channel
}
