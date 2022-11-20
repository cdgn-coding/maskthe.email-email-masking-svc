package queues

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
	queue   string
}

func (p *Publisher) Dispatch(message string) error {
	return p.channel.PublishWithContext(context.Background(),
		"",
		p.queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func NewPublisher(channel *amqp.Channel, name string) *Publisher {
	DeclareQueue(channel, name)
	return &Publisher{channel: channel, queue: name}
}
