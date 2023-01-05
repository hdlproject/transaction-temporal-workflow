package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	"time"
	"transaction-temporal-workflow/model"
)

type Transaction struct {
	redis *redis.Client
	rh    *rejson.Handler
}

func NewTransaction(redis *redis.Client) Transaction {
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(redis)

	return Transaction{
		redis: redis,
		rh:    rh,
	}
}

func (i Transaction) GetTransaction(transactionId string) (model.Transaction, error) {
	res, err := i.rh.JSONGet(transactionId, ".")
	if err != nil {
		return model.Transaction{}, err
	}

	var transaction model.Transaction
	err = json.Unmarshal(res.([]byte), &transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (i Transaction) CreateTransaction(transactionId string) error {
	transaction := model.Transaction{
		Id:     transactionId,
		Status: model.TransactionStatusCreated,
	}

	fmt.Println(transaction)

	_, err := i.rh.JSONSet(transactionId, ".", transaction)
	if err != nil {
		return err
	}

	return nil
}

func (i Transaction) UpdateTransactionStatus(transactionId string, transactionStatus model.TransactionStatus) error {
	_, err := i.rh.JSONSet(transactionId, "$.status", transactionStatus)
	if err != nil {
		return err
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
