FROM golang:1.20-buster AS build

# non-go modules dependencies
RUN apt update && apt install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

# copy code, protobuf and wait-for-it script
WORKDIR /go/src/app
COPY ./booking .
COPY ./wait-for-it.sh .

RUN go mod download
RUN go generate ./...
RUN go install
RUN go build

# separate production stage to reduce image size
FROM golang:1.20-buster AS production
WORKDIR /go/bin
COPY --from=build /go/src/app/booking .
COPY --from=build /go/src/app/docker-entrypoint.sh .
COPY --from=build /go/src/app/wait-for-it.sh .

RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["booking"]

EXPOSE 9112