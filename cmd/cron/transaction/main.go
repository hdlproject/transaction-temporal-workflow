package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/usecase"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(fmt.Errorf("new temporal client: %w", err))
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:           "transaction_relay",
		TaskQueue:    usecase.TransactionTaskQueue,
		CronSchedule: "* * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cmd.TransactionWorkflow.PublishTransaction)
	if err != nil {
		log.Fatalf("execute workflow: %v", err)
	}

	printResults(we.GetID(), we.GetRunID())
}

func printResults(workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
}
