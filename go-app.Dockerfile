FROM golang:1.18-alpine

WORKDIR /app

#RUN /bin/echo "::set-output name=go-build::$(go env GOCACHE)"
#RUN /bin/echo "::set-output name=go-mod::$(go env GOMODCACHE)"

COPY . ./
RUN go mod tidy

WORKDIR appdir
RUN go build -o appname

CMD [ "./appname" ]
