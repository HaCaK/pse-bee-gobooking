FROM golang:1.20-buster

# copy code
WORKDIR /go/src/app
COPY ./proxy .
COPY ./proto ./proto

RUN apt update && apt install -y protobuf-compiler
RUN GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

RUN go mod download
RUN go generate ./...
RUN go mod tidy -e
RUN go install

CMD ["proxy"]

EXPOSE 8080