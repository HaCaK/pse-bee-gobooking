package main

import (
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	"github.com/HaCaK/pse-bee-gobooking/src/property/service"
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

const grpcPort = 9111

func main() {
	log.Info("Starting goBooking property gRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen on grpc port %d: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterPropertyExternalServer(grpcServer, &service.PropertyService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
