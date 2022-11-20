package queues

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Invoke(string) error
}

func ConsumeQueue(channel *amqp.Channel, name string, consumer Consumer) {
	queue := DeclareQueue(channel, name)

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for d := range msgs {
			_ = consumer.Invoke(string(d.Body))
		}
	}()
}
