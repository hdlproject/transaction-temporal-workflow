package main

import (
	"context"
	"fmt"
	"os"

	"go.temporal.io/sdk/client"

	"temporalio-poc/api"
	"temporalio-poc/cmd"
	"temporalio-poc/model"
	"temporalio-poc/usecase"
)

type transactionServer struct {
	api.UnimplementedTransactionServer
	c client.Client
}

func NewTransactionServer() (api.TransactionServer, error) {
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TemporalAddress"),
	})
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
	we, err := s.c.ExecuteWorkflow(ctx, options, cmd.TransactionWorkflow.CreateTransaction, transactionReq)
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
	we, err := s.c.ExecuteWorkflow(ctx, options, cmd.TransactionWorkflow.ProcessTransaction, req.TransactionId, model.TransactionStatusSuccess)
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
