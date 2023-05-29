package integration_test

import (
	"context"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type MockPropertyInternalServer struct {
	proto.PropertyInternalServer
}

// Start creates and starts a mock PropertyInternalServer that listens on the given port
// // this is done to isolate testing of booking service from the actual implementation of the property service
func (h *MockPropertyInternalServer) Start(port string) func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on grpc port %s: %v", port, err)
	}

	baseServer := grpc.NewServer()
	proto.RegisterPropertyInternalServer(baseServer, h)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("Error serving propertyInternalServer: %v", err)
		}
	}()

	closer := func() {
		baseServer.Stop()
	}

	return closer
}

func (h *MockPropertyInternalServer) ConfirmBooking(_ context.Context, _ *proto.BookingReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}

func (h *MockPropertyInternalServer) CancelBooking(_ context.Context, _ *proto.BookingReq) (*emptypb.Empty, error) {
	return new(emptypb.Empty), nil
}
