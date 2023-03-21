import express from "express";

export class NotFoundError extends Error {
    constructor(message: string) {
        super(message);
        Object.setPrototypeOf(this, new.target.prototype);
        Error.captureStackTrace(this);
    }
}

export class ApiErrorHandler {
    handle(error: Error, response: express.Response) {
        if (error instanceof NotFoundError) {
            return response.status(401).send(error.message)
        } else {
            return response.status(500).send('internal server error')
        }
    }
}
