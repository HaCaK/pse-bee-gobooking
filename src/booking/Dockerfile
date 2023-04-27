FROM golang:1.20-buster

# copy code and protobuf
WORKDIR /go/src/app
COPY ./booking .

# non-go modules dependencies
RUN apt update && apt install -y protobuf-compiler
RUN GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

RUN go mod download
RUN go generate ./...
RUN go mod tidy -e
RUN go install

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/81b1373f17855a4dc21156cfe1694c31d7d1792e/wait-for-it.sh
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["booking"]

EXPOSE 9112