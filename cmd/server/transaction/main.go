package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"transaction-temporal-workflow/api"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9090))
	if err != nil {
		log.Fatalf("network listen: %v", err)
	}

	server, err := NewTransactionServer()
	if err != nil {
		log.Fatalf("new transaction server: %v", err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterTransactionServer(grpcServer, server)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("grpc server serve: %v", err)
	}
}

func printResults(workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
}
