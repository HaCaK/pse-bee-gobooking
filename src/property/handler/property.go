package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	"github.com/HaCaK/pse-bee-gobooking/src/property/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PropertyHandler struct {
	proto.UnimplementedPropertyExternalServer
	proto.UnimplementedPropertyInternalServer
}

func (h *PropertyHandler) CreateProperty(_ context.Context, req *proto.CreatePropertyReq) (*proto.Property, error) {
	property := model.Property{
		Name:        req.Name,
		Description: req.Description,
		OwnerName:   req.OwnerName,
		Address:     req.Address,
	}

	if err := service.CreateProperty(&property); err != nil {
		log.Errorf("Failure creating property: %v", err)
		return nil, err
	}

	return mapToProtoProperty(property), nil
}

func (h *PropertyHandler) UpdateProperty(_ context.Context, req *proto.UpdatePropertyReq) (*proto.Property, error) {
	property := model.Property{
		Name:        req.Name,
		Description: req.Description,
		OwnerName:   req.OwnerName,
		Address:     req.Address,
	}

	updatedProperty, err := service.UpdateProperty(uint(req.Id), &property)
	if err != nil {
		log.Errorf("Failure updating property with ID %v: %v", req.Id, err)
		return nil, err
	}
	if updatedProperty == nil {
		return nil, errors.New("404 property not found")
	}

	return mapToProtoProperty(*updatedProperty), nil
}

func (h *PropertyHandler) GetProperty(_ context.Context, req *proto.PropertyIdReq) (*proto.Property, error) {
	property, err := service.GetProperty(uint(req.Id))
	if err != nil {
		log.Errorf("Failure retrieving property with ID %v: %v", req.Id, err)
		return nil, err
	}
	if property == nil {
		return nil, errors.New("404 property not found")
	}
	return mapToProtoProperty(*property), nil
}

func (h *PropertyHandler) GetProperties(_ context.Context, _ *emptypb.Empty) (*proto.ListPropertiesResp, error) {
	properties, err := service.GetProperties()
	if err != nil {
		log.Errorf("Failure retrieving properties: %v", err)
		return nil, err
	}

	var protoProperties []*proto.Property
	for _, property := range properties {
		protoProperties = append(protoProperties, mapToProtoProperty(property))
	}
	return &proto.ListPropertiesResp{Properties: protoProperties}, nil
}

func (h *PropertyHandler) DeleteProperty(_ context.Context, req *proto.PropertyIdReq) (*emptypb.Empty, error) {
	property, err := service.DeleteProperty(uint(req.Id))
	if err != nil {
		log.Errorf("Failure deleting property with ID %v: %v", req.Id, err)
		return nil, err
	}
	if property == nil {
		return nil, errors.New("404 property not found")
	}
	return new(emptypb.Empty), nil
}

func (h *PropertyHandler) ConfirmBooking(_ context.Context, req *proto.BookingReq) (*emptypb.Empty, error) {
	log.Infof("Received booking request: %v", req)

	existingProperty, err := service.GetProperty(uint(req.PropertyId))
	if existingProperty == nil || err != nil {
		return nil, err
	}

	if existingProperty.IsStatusBooked() {
		return nil, fmt.Errorf("sorry, property %s is already booked", existingProperty.Name)
	}

	err = service.BookProperty(existingProperty, uint(req.BookingId))
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}

func (h *PropertyHandler) CancelBooking(_ context.Context, req *proto.BookingReq) (*emptypb.Empty, error) {
	log.Infof("Received cancellation request: %v", req)

	existingProperty, err := service.GetProperty(uint(req.PropertyId))
	if existingProperty == nil || err != nil {
		return nil, err
	}

	if existingProperty.BookingId != uint(req.BookingId) {
		return nil, errors.New("property is already booked by other booking")
	}

	err = service.FreeProperty(existingProperty)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}
