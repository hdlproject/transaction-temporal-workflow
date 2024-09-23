package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/usecase"
)

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TemporalAddress"),
	})
	if err != nil {
		panic(fmt.Errorf("new temporal client: %w", err))
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:           "user_relay",
		TaskQueue:    usecase.UserTaskQueue,
		CronSchedule: "* * * * *",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cmd.UserWorkflow.PublishUserBalanceEvent)
	if err != nil {
		log.Fatalf("execute workflow: %v", err)
	}

	printResults(we.GetID(), we.GetRunID())
}

func printResults(workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
}
