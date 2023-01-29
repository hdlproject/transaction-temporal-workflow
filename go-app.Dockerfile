FROM golang:1.18-alpine

WORKDIR /app

COPY . ./
RUN go mod tidy

#WORKDIR cmd/server/transaction
WORKDIR appdir
RUN go build -o appname

EXPOSE 8080

CMD [ "/appname" ]
