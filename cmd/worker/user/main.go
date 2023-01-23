package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"transaction-temporal-workflow/cmd"
	"transaction-temporal-workflow/dependency"
	"transaction-temporal-workflow/usecase"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	initRabbitMQ(cmd.RabbitMQ)

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, usecase.UserTaskQueue, worker.Options{})
	w.RegisterWorkflow(cmd.UserWorkflow.ReserveUserBalance)
	w.RegisterWorkflow(cmd.UserWorkflow.PublishUserBalanceEvent)

	w.RegisterActivity(cmd.UserActivity.ReserveUserBalance)
	w.RegisterActivity(cmd.UserActivity.PublishUserBalanceEvent)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func initRabbitMQ(rabbitMQ *amqp.Channel) {
	dependency.AddExchange(rabbitMQ, usecase.TransactionExchangeName)
}
