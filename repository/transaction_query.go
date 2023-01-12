package repository

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
)

type TransactionQuery struct {
	redis *redis.Client

	db *gorm.DB
}

func NewTransactionQuery(db *gorm.DB) TransactionQuery {
	return TransactionQuery{
		db: db,
	}
}

func (i TransactionQuery) GetTransactionByTransactionId(transactionId string) (transaction model.Transaction, err error) {
	result := i.db.Joins("Product").Joins("User").First(&transaction, "transaction_id = ?", transactionId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Transaction{}, result.Error
		}

		return model.Transaction{}, fmt.Errorf("get transaction by transaction id: %w", result.Error)
	}

	return transaction, nil
}

func (i TransactionQuery) GetLastTransactionByTransactionId(transactionId string) (transaction model.Transaction, err error) {
	result := i.db.Joins("Product").Joins("User").Order("created_at DESC").First(&transaction, "transaction_id = ?", transactionId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Transaction{}, result.Error
		}

		return model.Transaction{}, fmt.Errorf("get transaction by transaction id: %w", result.Error)
	}

	return transaction, nil
}
