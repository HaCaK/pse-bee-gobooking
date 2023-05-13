package main

import (
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/db"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/handler"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	"google.golang.org/grpc"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// ensure that logger is initialized before connecting to DB
	defer db.Init()
	// init logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
	level, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Info("Log level not specified, using default log level: INFO")
		log.SetLevel(log.InfoLevel)
		return
	}
	log.SetLevel(level)
}

var port = os.Getenv("PORT")

// main creates a gRPC server for all requests related to bookings
func main() {
	log.Info("Starting goBooking booking gRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	bookingHandler := new(handler.BookingHandler)
	proto.RegisterBookingExternalServer(grpcServer, bookingHandler)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
