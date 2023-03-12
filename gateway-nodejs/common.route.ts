import express from 'express'

export abstract class CommonRoute {
    app: express.Application
    name: string

    protected constructor(app: express.Application, name: string) {
        this.app = app
        this.name = name
        this.configureRoutes();
    }

    getName() {
        return this.name
    }

    abstract configureRoutes(): express.Application;
}