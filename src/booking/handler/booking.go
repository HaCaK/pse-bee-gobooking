package handler

import (
	"context"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/model"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BookingHandler struct {
	proto.BookingExternalServer
}

func (h *BookingHandler) CreateBooking(_ context.Context, req *proto.CreateBookingReq) (*proto.Booking, error) {
	booking := model.Booking{
		Comment:      req.Comment,
		CustomerName: req.CustomerName,
		PropertyId:   uint(req.PropertyId),
	}

	createdBooking, err := service.CreateBooking(&booking)
	if err != nil {
		log.Errorf("Error calling service CreateBooking: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return mapToProtoBooking(createdBooking), nil
}

func (h *BookingHandler) UpdateBooking(_ context.Context, req *proto.UpdateBookingReq) (*proto.Booking, error) {
	booking := model.Booking{
		Comment:      req.Comment,
		CustomerName: req.CustomerName,
		PropertyId:   uint(req.PropertyId),
	}

	updatedBooking, err := service.UpdateBooking(uint(req.Id), &booking)
	if err != nil {
		log.Errorf("Error calling service UpdateBooking with ID %v: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if updatedBooking == nil {
		return nil, status.Errorf(codes.NotFound, "Booking not found")
	}
	return mapToProtoBooking(updatedBooking), nil
}

func (h *BookingHandler) GetBooking(_ context.Context, req *proto.BookingIdReq) (*proto.Booking, error) {
	booking, err := service.GetBooking(uint(req.Id))
	if err != nil {
		log.Errorf("Error calling service GetBooking with ID %v: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if booking == nil {
		return nil, status.Errorf(codes.NotFound, "Booking not found")
	}
	return mapToProtoBooking(booking), nil
}

func (h *BookingHandler) GetBookings(_ context.Context, _ *emptypb.Empty) (*proto.ListBookingsResp, error) {
	bookings, err := service.GetBookings()
	if err != nil {
		log.Errorf("Error calling service GetBookings: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var protoBookings []*proto.Booking
	for _, booking := range bookings {
		protoBookings = append(protoBookings, mapToProtoBooking(&booking))
	}
	return &proto.ListBookingsResp{Bookings: protoBookings}, nil
}

func (h *BookingHandler) DeleteBooking(_ context.Context, req *proto.BookingIdReq) (*emptypb.Empty, error) {
	booking, err := service.DeleteBooking(uint(req.Id))
	if err != nil {
		log.Errorf("Error calling service DeleteBooking with ID %v: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if booking == nil {
		return nil, status.Errorf(codes.NotFound, "Booking not found")
	}
	return new(emptypb.Empty), nil
}
