package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"transaction-temporal-workflow/api"
	"transaction-temporal-workflow/usecase/idempotency"
)

type interceptor struct {
	idempotencyUseCase idempotency.UseCase
}

func NewInterceptor(idempotencyUseCase idempotency.UseCase) interceptor {
	return interceptor{
		idempotencyUseCase: idempotencyUseCase,
	}
}

func (i interceptor) checkIdempotency(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var idempotencyKey string
	switch v := req.(type) {
	case *api.CreateTransactionRequest:
		idempotencyKey = v.IdempotencyKey
	case *api.ProcessTransactionRequest:
		idempotencyKey = v.IdempotencyKey
	default:
		return nil, fmt.Errorf("request type is not supported")
	}

	isAllowed, err := i.idempotencyUseCase.IsAllowed(idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("is allowed: %w", err)
	}
	if !isAllowed {
		return nil, fmt.Errorf("idempotency key already exists")
	}

	return handler(ctx, req)
}
