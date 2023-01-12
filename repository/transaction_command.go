package repository

import (
	"fmt"

	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
)

type TransactionCommand struct {
	db *gorm.DB
}

func NewTransactionCommand(db *gorm.DB) TransactionCommand {
	return TransactionCommand{
		db: db,
	}
}

func (i TransactionCommand) CreateTransaction(transaction model.Transaction) error {
	result := i.db.Create(&transaction)
	if result.Error != nil {
		return fmt.Errorf("create transaction: %w", result.Error)
	}

	return nil
}
