/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Metadata } from "@grpc/grpc-js";
import { Observable } from "rxjs";

export const protobufPackage = "transaction_temporal_workflow.api";

export interface CreateTransactionRequest {
  transaction_id: string;
  amount: string;
  product_code: string;
  user_id: string;
  idempotency_key: string;
}

export interface CreateTransactionResponse {
  message: string;
}

export interface ProcessTransactionRequest {
  transaction_id: string;
  idempotency_key: string;
}

export interface ProcessTransactionResponse {
  message: string;
}

export const TRANSACTION_TEMPORAL_WORKFLOW_API_PACKAGE_NAME =
  "transaction_temporal_workflow.api";

export interface TransactionClient {
  createTransaction(
    request: CreateTransactionRequest,
    metadata?: Metadata
  ): Observable<CreateTransactionResponse>;

  processTransaction(
    request: ProcessTransactionRequest,
    metadata?: Metadata
  ): Observable<ProcessTransactionResponse>;
}

export interface TransactionController {
  createTransaction(
    request: CreateTransactionRequest,
    metadata?: Metadata
  ):
    | Promise<CreateTransactionResponse>
    | Observable<CreateTransactionResponse>
    | CreateTransactionResponse;

  processTransaction(
    request: ProcessTransactionRequest,
    metadata?: Metadata
  ):
    | Promise<ProcessTransactionResponse>
    | Observable<ProcessTransactionResponse>
    | ProcessTransactionResponse;
}

export function TransactionControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["createTransaction", "processTransaction"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(
        constructor.prototype,
        method
      );
      GrpcMethod("Transaction", method)(
        constructor.prototype[method],
        method,
        descriptor
      );
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(
        constructor.prototype,
        method
      );
      GrpcStreamMethod("Transaction", method)(
        constructor.prototype[method],
        method,
        descriptor
      );
    }
  };
}

export const TRANSACTION_SERVICE_NAME = "Transaction";
