package repository

import (
	"fmt"

	"gorm.io/gorm"

	"transaction-temporal-workflow/model"
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

func (i User) DeductUserBalance(userId string, amount int) error {
	result := i.db.Exec(`UPDATE "user" SET balance = ? WHERE id = ?`, gorm.Expr("balance - ?", amount), userId)
	if result.Error != nil {
		return fmt.Errorf("deduct user balance: %w", result.Error)
	}

	return nil
}
