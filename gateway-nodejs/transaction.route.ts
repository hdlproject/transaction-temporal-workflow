import express, {NextFunction} from 'express'
import {CommonRoute} from "./common.route";
import {TransactionUseCase} from "./transaction.usecase";
import {CreateTransactionRequest, ProcessTransactionRequest} from "./api/transaction";

export class TransactionRoute extends CommonRoute {
    transactionUseCase: TransactionUseCase

    constructor(app: express.Application, transactionUseCase: TransactionUseCase) {
        super(app, "Transaction")
        this.transactionUseCase = transactionUseCase
    }

    configureRoutes() {
        this.app.route(`/transaction`)
            .post(async (req: express.Request, res: express.Response, next: NextFunction) => {
                const request: CreateTransactionRequest = {
                    transactionId: req.body.transactionId,
                    amount: parseInt(req.body.amount),
                    productCode: req.body.productCode,
                    userId: req.body.userId,
                    idempotencyKey: req.body.idempotencyKey,
                }

                try {
                    let response = await this.transactionUseCase.createTransaction(request)
                    res.status(200).send(response.message);
                } catch (error) {
                    next(error)
                }
            });

        this.app.route(`/transaction/process`)
            .post(async (req: express.Request, res: express.Response, next: NextFunction) => {
                const request: ProcessTransactionRequest = {
                    transactionId: req.body.transactionId,
                    idempotencyKey: req.body.idempotencyKey,
                }

                try {
                    let response = await this.transactionUseCase.processTransaction(request)
                    res.status(200).send(response.message);
                } catch (error) {
                    next(error)
                }
            });

        return this.app;
    }
}