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
        const createTransactionPromise = (request: CreateTransactionRequest): Promise<CreateTransactionResponse> => {
            return new Promise((resolve: (response: CreateTransactionResponse) => any, reject: (err: Error) => void) => {
                this.transactionServerClient.createTransaction(request, (err: ServiceError | null, response: CreateTransactionResponse) => {
                    if (err) return reject(err)
                    resolve(response)
                })
            })
        }

        return createTransactionPromise(request)
    }

    processTransaction(request: ProcessTransactionRequest) {
        const processTransactionPromise = (request: ProcessTransactionRequest): Promise<ProcessTransactionResponse> => {
            return new Promise((resolve: (response: ProcessTransactionResponse) => any, reject: (err: Error) => void) => {
                this.transactionServerClient.processTransaction(request, (err: ServiceError | null, response: ProcessTransactionResponse) => {
                    if (err) return reject(err)
                    resolve(response)
                })
            })
        }

        return processTransactionPromise(request)
    }
}
