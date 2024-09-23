package repository

import (
	"fmt"

	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"

	"temporalio-poc/model"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	return User{
		db: db,
	}
}

func (i User) GetUserById(id string) (user model.User, err error) {
	result := i.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return model.User{}, fmt.Errorf("get user by id: %w", result.Error)
	}

	return user, nil
}

func (i User) DeductUserBalance(transaction model.Transaction) error {
	totalPrice, err := transaction.GetTotalPrice()
	if err != nil {
		return fmt.Errorf("get total price: %w", err)
	}

	err = i.db.Transaction(func(tx *gorm.DB) error {
		transactionStatus := model.TransactionStatusPending
		result := tx.Exec(`UPDATE "user" SET balance = ? WHERE id = ?`, gorm.Expr("balance - ?", totalPrice), transaction.UserId)
		if result.Error != nil {
			transactionStatus = model.TransactionStatusFailed
			log.Error(fmt.Errorf("deduct user balance: %w", result.Error))
		}

		userBalanceEvent := model.UserBalanceEvent{
			UserId:            transaction.UserId,
			Balance:           -1 * totalPrice,
			TransactionId:     transaction.TransactionId,
			TransactionStatus: transactionStatus,
			IsPublished:       false,
		}

		if err := tx.Create(&userBalanceEvent).Error; err != nil {
			return fmt.Errorf("create user balance event : %w", err)
		}

		return nil
	})
	return err
}
