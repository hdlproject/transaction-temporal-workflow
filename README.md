## How to
### Run infrastructure
```shell
$ docker-compose up
```
### Run worker
```shell
$ go run cmd/worker/main.go
```
### Run server
```shell
$ go run cmd/server/*
```
### Open temporal admin UI
Open http://localhost:8080 in browser
### Open grpcui
```shell
$ grpcui -plaintext -port=9091 localhost:9090
```
Open http://localhost:9091 in browser
