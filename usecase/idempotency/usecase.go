package idempotency

import (
	"fmt"

	"transaction-temporal-workflow/repository"
)

type (
	UseCase struct {
		idempotencyRepository repository.Idempotency
	}
)

func NewUseCase(idempotencyRepository repository.Idempotency) UseCase {
	return UseCase{
		idempotencyRepository: idempotencyRepository,
	}
}

func (i UseCase) IsAllowed(idempotencyKey string) (bool, error) {
	isAllowed, err := i.idempotencyRepository.IsAllowed(idempotencyKey)
	if err != nil {
		return false, fmt.Errorf("is allowed: %w", err)
	}

	return isAllowed, nil
}
