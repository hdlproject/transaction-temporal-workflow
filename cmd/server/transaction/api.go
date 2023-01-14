package main

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/api"
	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/usecase"
)

type transactionServer struct {
	api.UnimplementedTransactionServer
	c client.Client
}

func NewTransactionServer() (api.TransactionServer, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("new temporal client: %w", err)
	}

	return &transactionServer{
		c: c,
	}, nil
}

func (s *transactionServer) CreateTransaction(ctx context.Context, req *api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	transactionReq := model.Transaction{
		TransactionId: req.TransactionId,
		Amount:        int(req.Amount),
		ProductCode:   req.ProductCode,
		UserId:        req.UserId,
	}

	options := client.StartWorkflowOptions{
		ID:        "transaction-workflow",
		TaskQueue: usecase.TransactionTaskQueue,
	}
	we, err := s.c.ExecuteWorkflow(ctx, options, cmd.TransactionWorkflow.CreateTransaction, transactionReq, req.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("execute workflow: %w", err)
	}

	err = we.Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("get workflow result: %w", err)
	}

	printResults(we.GetID(), we.GetRunID())

	return &api.CreateTransactionResponse{
		Message: "success",
	}, nil
}

func (s *transactionServer) ProcessTransaction(ctx context.Context, req *api.ProcessTransactionRequest) (*api.ProcessTransactionResponse, error) {
	options := client.StartWorkflowOptions{
		ID:        "transaction-workflow",
		TaskQueue: usecase.TransactionTaskQueue,
	}
	we, err := s.c.ExecuteWorkflow(ctx, options, cmd.TransactionWorkflow.ProcessTransaction, req.TransactionId, req.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("execute workflow: %w", err)
	}

	err = we.Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("get workflow result: %w", err)
	}

	printResults(we.GetID(), we.GetRunID())

	return &api.ProcessTransactionResponse{
		Message: "success",
	}, nil
}
