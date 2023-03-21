import express from 'express';
import * as http from 'http';
import * as winston from 'winston';
import * as expressWinston from 'express-winston';
import cors from 'cors';
import {CommonRoute} from './common.route';
import {TransactionRoute} from './transaction.route';
import {TransactionUseCase} from './transaction.usecase';
import debug from 'debug';
import {AdminRoute} from "./admin.route";
import {ApiErrorHandler} from "./api.error";

const app: express.Application = express();
app.use(express.json());
app.use(cors());

const loggerOptions: expressWinston.LoggerOptions = {
    transports: [new winston.transports.Console()],
    format: winston.format.combine(
        winston.format.json(),
        winston.format.prettyPrint(),
        winston.format.colorize({all: true})
    ),
};
if (!process.env.DEBUG) {
    loggerOptions.meta = false;
}
app.use(expressWinston.logger(loggerOptions));

const server: http.Server = http.createServer(app);
const port = 3000;

const routes: Array<CommonRoute> = [];
routes.push(new TransactionRoute(app, new TransactionUseCase()))
routes.push(new AdminRoute(app, port))

const apiErrorHandler = new ApiErrorHandler();
app.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
    return apiErrorHandler.handle(err, res)
})

const debugLog: debug.IDebugger = debug('app');
server.listen(port, () => {
    routes.forEach((route: CommonRoute) => {
        debugLog(`Routes configured for ${route.getName()}`);
    });
});
