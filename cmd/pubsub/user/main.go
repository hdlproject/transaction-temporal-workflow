package main

import (
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/cmd"
)

func main() {
	msgs, err := cmd.RabbitMQ.Consume(
		cmd.UserServiceQueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("rabbitmq consume: %v", err))
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			options := client.StartWorkflowOptions{
				ID:        "transaction-workflow",
				TaskQueue: cmd.UserTaskQueue,
			}
			we, err := s.c.ExecuteWorkflow(ctx, options, cmd.TransactionWorkflow.CreateTransaction, transactionReq, req.IdempotencyKey)
			if err != nil {
				return nil, fmt.Errorf("execute workflow: %w", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func printResults(workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
}
