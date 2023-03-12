import express from 'express'
import {CommonRoute} from "./common.route";

export class AdminRoute extends CommonRoute {
    port: number

    constructor(app: express.Application, port: number) {
        super(app, "Admin")
        this.port = port
    }

    configureRoutes() {
        this.app.get('/admin/health-check', (req: express.Request, res: express.Response) => {
            const runningMessage = `Server running at http://localhost:${this.port}`;
            res.status(200).send(runningMessage)
        });

        return this.app;
    }
}
