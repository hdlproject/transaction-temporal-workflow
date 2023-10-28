package dependency

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ() *amqp.Channel {
	conn, err := amqp.Dial(os.Getenv("RabbitMQURL"))
	if err != nil {
		panic("amqp dial: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("conn channel: " + err.Error())
	}

	return ch
}

func AddExchange(ch *amqp.Channel, name string) {
	err := ch.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("exchange declare: " + err.Error())
	}
}

func AddQueue(ch *amqp.Channel, name string) {
	_, err := ch.QueueDeclare(
		name,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		panic("queue declare: " + err.Error())
	}
}

func AddRouting(ch *amqp.Channel, exchangeName, name, routingKey string) {
	err := ch.QueueBind(
		name,
		routingKey,
		exchangeName,
		false,
		nil)
	if err != nil {
		panic("queue bind: " + err.Error())
	}
}
