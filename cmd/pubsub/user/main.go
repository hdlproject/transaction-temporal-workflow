package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/dependency"
	"transaction-temporal-workflow/model"
	"transaction-temporal-workflow/usecase"
)

func main() {
	initRabbitMQ(cmd.RabbitMQ)

	msgs, err := cmd.RabbitMQ.Consume(
		usecase.UserServiceQueueName,
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

	c, err := client.NewClient(client.Options{
		HostPort: os.Getenv("TemporalAddress"),
	})
	if err != nil {
		panic(fmt.Errorf("new temporal client: %w", err))
	}
	defer c.Close()

	var forever chan struct{}

	go func(c client.Client) {
		for d := range msgs {
			ctx := context.Background()

			options := client.StartWorkflowOptions{
				ID:        "user-workflow",
				TaskQueue: usecase.UserTaskQueue,
			}
			var we client.WorkflowRun
			switch d.RoutingKey {
			case usecase.TransactionCreatedRoutingKey:
				var transaction model.Transaction
				err := json.Unmarshal(d.Body, &transaction)
				if err != nil {
					log.Printf("json unmarshal: %v", err)
					return
				}

				idempotencyKey := fmt.Sprintf("%s.%s", usecase.TransactionCreatedRoutingKey, transaction.TransactionId)
				isAllowed, err := cmd.IdempotencyUseCase.IsAllowed(idempotencyKey)
				if err != nil {
					log.Printf("is allowed: %v", err)
					return
				}
				if !isAllowed {
					log.Printf("message with idempotency key %s is not allowed", idempotencyKey)
					continue
				}

				we, err = c.ExecuteWorkflow(ctx, options, cmd.UserWorkflow.ReserveUserBalance, transaction)
				if err != nil {
					log.Printf("execute workflow: %v", err)
					return
				}
			default:
				log.Printf("routing key %s is not supported", d.RoutingKey)
				return
			}

			err = we.Get(context.Background(), nil)
			if err != nil {
				log.Printf("get workflow result: %v", err)
			}

			printResults(we.GetID(), we.GetRunID())
		}
	}(c)

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func printResults(workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
}

func initRabbitMQ(rabbitMQ *amqp.Channel) {
	dependency.AddExchange(rabbitMQ, usecase.TransactionExchangeName)
	dependency.AddQueue(rabbitMQ, usecase.UserServiceQueueName)
	dependency.AddRouting(rabbitMQ, usecase.TransactionExchangeName, usecase.UserServiceQueueName, usecase.TransactionCreatedRoutingKey)
}
