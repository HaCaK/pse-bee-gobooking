package handler

import (
	"context"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/db"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/model"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

// creates and starts a BookingExternalServer and returns a client that is connected to it and can be used for tests
func startBookingExternalServer(ctx context.Context) (proto.BookingExternalClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	proto.RegisterBookingExternalServer(baseServer, new(BookingHandler))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("Error serving bookingExternalServer: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error connecting to bookingExternalServer: %v", err)
	}

	closer := func() {
		baseServer.Stop()
	}

	client := proto.NewBookingExternalClient(conn)

	return client, closer
}

func getMockListBookingsResp(bookingResp *proto.BookingResp) *proto.ListBookingsResp {
	list := new(proto.ListBookingsResp)

	if bookingResp != nil {
		list.Bookings = append(list.Bookings, bookingResp)
	}

	return list
}

func getMockBookingRespWithDefaultCustomerName() *proto.BookingResp {
	return getMockBookingResp("customer")
}

func getMockBookingResp(customerName string) *proto.BookingResp {
	return &proto.BookingResp{
		Id:           1,
		Comment:      "comment",
		CustomerName: customerName,
		Status:       "PENDING",
		PropertyId:   1,
	}
}

func createBookingInDB() {
	booking := model.Booking{
		Comment:      "comment",
		CustomerName: "customer",
		Status:       "PENDING",
		PropertyId:   1,
	}
	db.DB.Create(&booking)
}

func deleteBookingInDB() {
	db.DB.Delete(new(model.Booking), 1)
}
