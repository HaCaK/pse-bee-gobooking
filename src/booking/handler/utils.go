package handler

import (
	"github.com/HaCaK/pse-bee-gobooking/src/booking/model"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapToProtoBooking(booking *model.Booking) *proto.Booking {
	return &proto.Booking{
		Id:           uint32(booking.ID),
		Comment:      booking.Comment,
		CustomerName: booking.CustomerName,
		Status:       string(booking.Status),
		PropertyId:   uint32(booking.PropertyId),
		CreatedAt:    timestamppb.New(booking.CreatedAt),
		UpdatedAt:    timestamppb.New(booking.UpdatedAt),
	}
}
