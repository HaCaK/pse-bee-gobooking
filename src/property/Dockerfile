FROM golang:1.20-buster

# copy code
WORKDIR /go/src/app
COPY ./property .

RUN apt update && apt install -y protobuf-compiler
RUN GO111MODULE=off go get google.golang.org/protobuf/cmd/protoc-property-go google.golang.org/grpc/cmd/protoc-property-go-grpc

RUN go mod download
RUN buf generate
RUN go mod tidy -e
RUN go install

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/81b1373f17855a4dc21156cfe1694c31d7d1792e/wait-for-it.sh
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["property"]

EXPOSE 9111