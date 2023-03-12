import {
    CreateTransactionRequest,
    CreateTransactionResponse,
    ProcessTransactionRequest,
    ProcessTransactionResponse,
    TransactionClient
} from './api/transaction';
import {credentials, ServiceError} from '@grpc/grpc-js';

export class TransactionUseCase {
    transactionServerClient: TransactionClient

    constructor() {
        this.transactionServerClient = new TransactionClient('localhost:9090', credentials.createInsecure());
    }

    createTransaction(request: CreateTransactionRequest) {
        this.transactionServerClient.createTransaction(request, (err: ServiceError | null, response: CreateTransactionResponse) => {
            console.log(response, err)
        });
    }

    processTransaction(request: ProcessTransactionRequest) {
        this.transactionServerClient.processTransaction(request, (err: ServiceError | null, response: ProcessTransactionResponse) => {
            console.log(response, err)
        });
    }
}
