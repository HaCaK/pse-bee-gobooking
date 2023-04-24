package main

import (
	"context"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/proto"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

const proxyPort = 8080

var propertyTarget = os.Getenv("PROPERTY_CONNECT")

func main() {
	mux := runtime.NewServeMux()
	err := proto.RegisterPropertyExternalHandlerFromEndpoint(context.Background(), mux, propertyTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	server := gin.New()
	server.Use(gin.Logger())
	server.Group("properties/*{grpc_gateway}").Any("", gin.WrapH(mux))

	// start server
	err = server.Run(fmt.Sprintf(":%d", proxyPort))
	if err != nil {
		log.Fatal(err)
	}
}
