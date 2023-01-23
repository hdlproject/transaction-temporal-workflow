## How to
### Run infrastructure
```shell
$ docker-compose up
```
### Transaction service
#### Run server
```shell
$ go run cmd/server/transaction/*
```
#### Run message listener
```shell
$ go run cmd/pubsub/transaction/*
```
#### Run cron
```shell
$ go run cmd/cron/transaction/*
```
#### Run worker
```shell
$ go run cmd/worker/transaction/*
```
### User service
#### Run message listener
```shell
$ go run cmd/pubsub/user/*
```
#### Run cron
```shell
$ go run cmd/cron/user/*
```
#### Run worker
```shell
$ go run cmd/worker/user/*
```
### Open temporal admin UI
Open http://localhost:8080 in browser
### Open grpcui
```shell
$ grpcui -plaintext -port=9091 localhost:9090
```
Open http://localhost:9091 in browser
