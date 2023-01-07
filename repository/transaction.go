package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
)

type Transaction struct {
	redis *redis.Client
	rh    *rejson.Handler

	db *gorm.DB
}

func NewTransaction(redis *redis.Client, db *gorm.DB) Transaction {
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(redis)

	return Transaction{
		redis: redis,
		rh:    rh,
		db:    db,
	}
}

func (i Transaction) GetTransactionByTransactionId(transactionId string) (transaction model.Transaction, err error) {
	result := i.db.First(&transaction, "transaction_id = ?", transactionId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Transaction{}, result.Error
		}

		return model.Transaction{}, fmt.Errorf("get transaction by transaction id: %w", result.Error)
	}

	return transaction, nil
}

func (i Transaction) GetLastTransactionByTransactionId(transactionId string) (transaction model.Transaction, err error) {
	result := i.db.Order("created_at DESC").First(&transaction, "transaction_id = ?", transactionId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.Transaction{}, result.Error
		}

		return model.Transaction{}, fmt.Errorf("get transaction by transaction id: %w", result.Error)
	}

	return transaction, nil
}

func (i Transaction) CreateTransaction(transaction model.Transaction) error {
	result := i.db.Create(&transaction)
	if result.Error != nil {
		return fmt.Errorf("create transaction: %w", result.Error)
	}

	return nil
}

func (i Transaction) IsAllowed(idempotencyKey string) (bool, error) {
	res := i.redis.SetNX(context.Background(), idempotencyKey, "", time.Hour)
	if res.Err() != nil {
		return false, res.Err()
	}

	return res.Val(), nil
}
