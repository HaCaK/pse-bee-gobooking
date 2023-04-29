package main

import (
	"context"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/proxy/proto"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var port = os.Getenv("PORT")

var propertyTarget = os.Getenv("PROPERTY_CONNECT")
var bookingTarget = os.Getenv("BOOKING_CONNECT")

// main creates a gRPC gateway which acts as a proxy between external HTTP clients
// and the internal gRPC property and booking services
func main() {
	// Register gRPC handlers for property and booking services
	mux := runtime.NewServeMux()
	err := proto.RegisterPropertyExternalHandlerFromEndpoint(context.Background(), mux, propertyTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	err = proto.RegisterBookingExternalHandlerFromEndpoint(context.Background(), mux, bookingTarget, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("Failed to connect to gRPC clients: %v", err)
	}

	// Create an HTTP server
	server := gin.New()
	server.Use(gin.Logger())

	handlerFunc := gin.WrapH(mux)

	// Forward any HTTP requests to gRPC handlers
	server.Group("properties").Any("", handlerFunc)
	server.Group("properties/*{grpc_gateway}").Any("", handlerFunc)

	server.Group("bookings").Any("", handlerFunc)
	server.Group("bookings/*{grpc_gateway}").Any("", handlerFunc)

	log.Info("Starting goBooking proxy server")
	err = server.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
