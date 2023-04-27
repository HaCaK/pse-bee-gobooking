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
	if err != nil {
		log.Fatal(err)
	}

	mux2 := runtime.NewServeMux()
	err = proto.RegisterBookingExternalHandlerFromEndpoint(context.Background(), mux2, bookingTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatal(err)
	}

	// Creating a normal HTTP server
	server := gin.New()
	server.Use(gin.Logger())

	propertyHandlerFunc := gin.WrapH(mux)
	bookingHandlerFunc := gin.WrapH(mux2)

	server.Group("properties").Any("", propertyHandlerFunc)
	server.Group("properties/*{grpc_gateway}").Any("", propertyHandlerFunc)

	server.Group("bookings").Any("", bookingHandlerFunc)
	server.Group("bookings/*{grpc_gateway}").Any("", bookingHandlerFunc)

	// start server
	err = server.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
