package repository

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Transaction struct {
	db       *gorm.DB
	rabbitMQ *amqp.Channel

	Command TransactionCommand
	Query   TransactionQuery
}

func NewTransaction(db *gorm.DB, rabbitMQ *amqp.Channel) Transaction {
	return Transaction{
		db:       db,
		rabbitMQ: rabbitMQ,
		Command:  NewTransactionCommand(db),
		Query:    NewTransactionQuery(db),
	}
}
