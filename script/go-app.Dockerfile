FROM golang:1.18-alpine

WORKDIR /app

COPY .. ./
RUN go mod tidy

WORKDIR appdir
RUN go build -o appname

CMD [ "./appname" ]
