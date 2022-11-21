package queues

import amqp "github.com/rabbitmq/amqp091-go"

func DeclareQueue(channel *amqp.Channel, name string) amqp.Queue {
	q, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		panic(err)
	}

	return q
}
