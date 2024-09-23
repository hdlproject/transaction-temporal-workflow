FROM golang:1.19.3-buster

WORKDIR /generate

RUN apt update \
    && apt install -y protobuf-compiler \
    && apt install -y curl \
    && curl -sL https://deb.nodesource.com/setup_14.x | bash -  \
    && apt install -y nodejs

COPY ../go.mod ../go.sum ./

RUN go install \
    github.com/golang/protobuf/protoc-gen-go

RUN npm install grpc-tools ts-proto@1.122.0 @bufbuild/protoc-gen-es@1.0.0 @bufbuild/protoc-gen-connect-es@0.8.0 -g --unsafe-perm --target_arch=x64

ENTRYPOINT ["bash"]
