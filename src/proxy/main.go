package main

import (
	"context"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/proxy/proto"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

var port = os.Getenv("PORT")

var propertyTarget = os.Getenv("PROPERTY_CONNECT")
var bookingTarget = os.Getenv("BOOKING_CONNECT")

func main() {
	mux := runtime.NewServeMux()
	err := proto.RegisterPropertyExternalHandlerFromEndpoint(context.Background(), mux, propertyTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	err = proto.RegisterBookingExternalHandlerFromEndpoint(context.Background(), mux, bookingTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	server := gin.New()
	server.Use(gin.Logger())

	handlerFunc := gin.WrapH(mux)

	server.Group("properties").Any("", handlerFunc)
	server.Group("properties/*{grpc_gateway}").Any("", handlerFunc)

	server.Group("bookings").Any("", handlerFunc)
	server.Group("bookings/*{grpc_gateway}").Any("", handlerFunc)

	// start server
	err = server.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
