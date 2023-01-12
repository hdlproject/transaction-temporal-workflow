package repository

import (
	"gorm.io/gorm"
)

type Transaction struct {
	db *gorm.DB

	Command TransactionCommand
	Query   TransactionQuery
}

func NewTransaction(db *gorm.DB) Transaction {
	return Transaction{
		db:      db,
		Command: NewTransactionCommand(db),
		Query:   NewTransactionQuery(db),
	}
}
