FROM golang:1.19.3-buster

WORKDIR /generate

RUN apt update \
    && apt install -y protobuf-compiler

COPY go.mod go.sum ./

RUN go install \
    github.com/golang/protobuf/protoc-gen-go

ENTRYPOINT ["bash"]
