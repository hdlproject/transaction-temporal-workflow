package repository

import (
	"fmt"

	"gorm.io/gorm"

	"temporalio-poc/model"
)

func (i User) GetUnpublishedUserBalanceEvents() (userBalanceEvents []model.UserBalanceEvent, err error) {
	result := i.db.Joins("User").Order("created_at ASC").Limit(10).Find(&userBalanceEvents, "is_published = FALSE")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}

		return nil, fmt.Errorf("get unpublished user balance event: %w", result.Error)
	}

	return userBalanceEvents, nil
}

func (i User) PublishUserBalanceEvent(id int64) error {
	result := i.db.Exec(`UPDATE user_balance_event SET is_published = TRUE WHERE id = ?`, id)
	if result.Error != nil {
		return fmt.Errorf("publish user balance event: %w", result.Error)
	}

	return nil
}
