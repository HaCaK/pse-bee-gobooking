FROM golang:1.20-buster AS build

# non-go modules dependencies
RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2

# copy code and external proto files
WORKDIR /go/src/app
COPY ./proxy .
COPY ./property/proto/property_external.proto ./proto/
COPY ./booking/proto/booking_external.proto ./proto/

RUN go mod download
RUN go generate ./...
RUN go install
RUN go build

# separate production stage to reduce image size
FROM golang:1.20-buster AS production
WORKDIR /go/bin
COPY --from=build /go/src/app/proxy .

CMD ["proxy"]

EXPOSE 8080