package handler

import (
	"context"
	"errors"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	"github.com/HaCaK/pse-bee-gobooking/src/property/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PropertyHandler struct {
	proto.PropertyExternalServer
	proto.PropertyInternalServer
}

func (h *PropertyHandler) CreateProperty(_ context.Context, req *proto.CreatePropertyReq) (*proto.PropertyResp, error) {
	property := model.Property{
		Name:        req.Name,
		Description: req.Description,
		OwnerName:   req.OwnerName,
		Address:     req.Address,
	}

	if err := service.CreateProperty(&property); err != nil {
		log.Errorf("Error calling service CreateProperty: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return mapToProtoPropertyResp(&property), nil
}

func (h *PropertyHandler) UpdateProperty(_ context.Context, req *proto.UpdatePropertyReq) (*proto.PropertyResp, error) {
	property := model.Property{
		Name:        req.Name,
		Description: req.Description,
		OwnerName:   req.OwnerName,
		Address:     req.Address,
	}

	updatedProperty, err := service.UpdateProperty(uint(req.Id), &property)
	if err != nil {
		log.Errorf("Error calling service UpdateProperty with ID %v: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if updatedProperty == nil {
		return nil, status.Errorf(codes.NotFound, "Property not found")
	}

	return mapToProtoPropertyResp(updatedProperty), nil
}

func (h *PropertyHandler) GetProperty(_ context.Context, req *proto.PropertyIdReq) (*proto.PropertyResp, error) {
	property, err := service.GetProperty(uint(req.Id))
	if err != nil {
		log.Errorf("Error calling service GetProperty with ID %v: %v", req.Id, err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if property == nil {
		return nil, status.Errorf(codes.NotFound, "Property not found")
	}
	return mapToProtoPropertyResp(property), nil
}

func (h *PropertyHandler) GetProperties(_ context.Context, _ *emptypb.Empty) (*proto.ListPropertiesResp, error) {
	properties, err := service.GetProperties()
	if err != nil {
		log.Errorf("Error calling service GetProperties: %v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var protoProperties []*proto.PropertyResp
	for _, property := range properties {
		protoProperties = append(protoProperties, mapToProtoPropertyResp(&property))
	}
	return &proto.ListPropertiesResp{Properties: protoProperties}, nil
}

func (h *PropertyHandler) DeleteProperty(_ context.Context, req *proto.PropertyIdReq) (*emptypb.Empty, error) {
	property, err := service.DeleteProperty(uint(req.Id))

	if err != nil {
		log.Errorf("Error calling service DeleteProperty with ID %v: %v", req.Id, err)

		var propertyError *model.PropertyError
		if errors.As(err, &propertyError) {
			return nil, status.Errorf(codes.InvalidArgument, propertyError.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if property == nil {
		return nil, status.Errorf(codes.NotFound, "Property not found")
	}
	return new(emptypb.Empty), nil
}

func (h *PropertyHandler) ConfirmBooking(_ context.Context, req *proto.BookingReq) (*emptypb.Empty, error) {
	log.Infof("Received booking request: %v", req)

	existingProperty, err := service.GetProperty(uint(req.PropertyId))
	if existingProperty == nil {
		return nil, status.Errorf(codes.NotFound, "Property not found")
	}
	if err != nil {
		return nil, err
	}

	err = service.BookProperty(existingProperty, uint(req.BookingId))
	if err != nil {
		log.Errorf("Error calling service BookProperty with ID %v: %v", req.PropertyId, err)

		var propertyError *model.PropertyError
		if errors.As(err, &propertyError) {
			return nil, status.Errorf(codes.InvalidArgument, propertyError.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return new(emptypb.Empty), nil
}

func (h *PropertyHandler) CancelBooking(_ context.Context, req *proto.BookingReq) (*emptypb.Empty, error) {
	log.Infof("Received cancellation request: %v", req)

	existingProperty, err := service.GetProperty(uint(req.PropertyId))
	if existingProperty == nil {
		return nil, status.Errorf(codes.NotFound, "Property not found")
	}
	if err != nil {
		return nil, err
	}

	err = service.FreeProperty(existingProperty, uint(req.BookingId))
	if err != nil {
		log.Errorf("Error calling service FreeProperty with ID %v: %v", req.PropertyId, err)

		var propertyError *model.PropertyError
		if errors.As(err, &propertyError) {
			return nil, status.Errorf(codes.InvalidArgument, propertyError.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return new(emptypb.Empty), nil
}
