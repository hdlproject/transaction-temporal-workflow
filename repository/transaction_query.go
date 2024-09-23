package repository

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"temporalio-poc/model"
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

func (i TransactionQuery) GetUnpublishedTransactions() (transactions []model.Transaction, err error) {
	result := i.db.Joins("Product").Joins("User").Order("created_at ASC").Limit(10).Find(&transactions, "is_published = FALSE")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}

		return nil, fmt.Errorf("get unpublished transactions: %w", result.Error)
	}

	return transactions, nil
}
