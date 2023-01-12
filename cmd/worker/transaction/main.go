package main

import (
	"log"

	"transaction-temporal-workflow/cmd"

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

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, cmd.TransactionTaskQueue, worker.Options{})
	w.RegisterWorkflow(cmd.TransactionWorkflow.CreateTransaction)
	w.RegisterWorkflow(cmd.TransactionWorkflow.ProcessTransaction)

	w.RegisterActivity(cmd.TransactionActivity.CreateTransaction)
	w.RegisterActivity(cmd.TransactionActivity.ProcessTransaction)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
