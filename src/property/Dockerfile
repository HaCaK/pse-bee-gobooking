FROM golang:1.20-buster AS build

# non-go modules dependencies
RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# copy code and protobuf
WORKDIR /go/src/app
COPY ./property .

RUN go mod download
RUN go generate ./...
RUN go install
RUN go build

# separate production stage to reduce image size
FROM golang:1.20-buster AS production
WORKDIR /go/bin
COPY --from=build /go/src/app/property .
COPY --from=build /go/src/app/docker-entrypoint.sh .
COPY ./wait-for-it.sh .

RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["property"]

EXPOSE 9111