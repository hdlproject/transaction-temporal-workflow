package main

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.temporal.io/sdk/client"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/dependency"
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

	c, err := client.NewClient(client.Options{})
	if err != nil {
		panic(fmt.Errorf("new temporal client: %w", err))
	}

	var forever chan struct{}

	go func(c client.Client) {
		for d := range msgs {
			ctx := context.Background()

			options := client.StartWorkflowOptions{
				ID:        "transaction-workflow",
				TaskQueue: usecase.UserTaskQueue,
			}
			we, err := c.ExecuteWorkflow(ctx, options, cmd.UserWorkflow.ConsumeMessage, d)
			if err != nil {
				log.Printf("execute workflow: %v", err)
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
