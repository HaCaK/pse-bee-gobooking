package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PropertyService struct {
	proto.UnimplementedPropertyExternalServer
	proto.UnimplementedPropertyInternalServer
}

func (s *PropertyService) CreateProperty(_ context.Context, req *proto.CreatePropertyReq) (*proto.Property, error) {
	property := model.Property{
		Name:        req.Name,
		Description: req.Description,
		OwnerName:   req.OwnerName,
		Address:     req.Address,
	}

	property.SetStatusFree()

	result := db.DB.Create(&property)
	if result.Error != nil {
		return nil, result.Error
	}

	log.Infof("Successfully stored new property with ID %v in database.", property.ID)
	log.Tracef("Stored: %v", property)
	return mapToProtoProperty(property), nil
}

func (s *PropertyService) UpdateProperty(_ context.Context, req *proto.UpdatePropertyReq) (*proto.Property, error) {
	existingProperty, err := getProperty(uint(req.Id))
	if existingProperty == nil || err != nil {
		return nil, err
	}

	existingProperty.Name = req.Name
	existingProperty.Description = req.Description
	existingProperty.OwnerName = req.OwnerName
	existingProperty.Address = req.Address
	result := db.DB.Save(existingProperty)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", req.Id)
	entry.Info("Successfully updated property.")
	entry.Tracef("Updated: %v", existingProperty)
	return mapToProtoProperty(*existingProperty), nil
}

func (s *PropertyService) ShowProperty(_ context.Context, req *proto.PropertyIdReq) (*proto.Property, error) {
	existingProperty, err := getProperty(uint(req.Id))
	if existingProperty == nil || err != nil {
		return nil, err
	}
	return mapToProtoProperty(*existingProperty), nil
}

func (s *PropertyService) ListProperties(_ context.Context, _ *emptypb.Empty) (*proto.ListPropertiesResp, error) {
	var properties []model.Property
	result := db.DB.Find(&properties)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", properties)

	var protoProperties []*proto.Property
	for _, property := range properties {
		protoProperties = append(protoProperties, mapToProtoProperty(property))
	}
	return &proto.ListPropertiesResp{Properties: protoProperties}, nil
}

func (s *PropertyService) DeleteProperty(_ context.Context, req *proto.PropertyIdReq) (*emptypb.Empty, error) {
	existingProperty, err := getProperty(uint(req.Id))
	if existingProperty == nil || err != nil {
		return new(emptypb.Empty), err
	}
	result := db.DB.Delete(existingProperty)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", req.Id)
	entry.Info("Successfully deleted property.")
	entry.Tracef("Deleted: %v", existingProperty)
	return new(emptypb.Empty), nil
}

// should be private because it is an internal func
func bookProperty(id uint, bookingId uint) (*model.Property, error) {
	existingProperty, err := getProperty(id)
	if existingProperty == nil || err != nil {
		return existingProperty, err
	}
	existingProperty.SetStatusBooked()
	existingProperty.BookingId = bookingId

	result := db.DB.Save(existingProperty)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", id)
	entry.Info("Successfully booked property.")
	entry.Tracef("Updated: %v", existingProperty)
	return existingProperty, nil
}

func (s *PropertyService) ProcessBooking(_ context.Context, booking *proto.BookingReq) error {
	log.Infof("Received booking: %v", booking)

	existingProperty, err := getProperty(uint(booking.PropertyId))
	if existingProperty == nil || err != nil {
		return err
	}

	if existingProperty.IsStatusBooked() {
		return fmt.Errorf("sorry, property %s is already booked", existingProperty.Name)
	}

	existingProperty.SetStatusBooked()

	_, err = bookProperty(existingProperty.ID, uint(booking.BookingId))
	if err != nil {
		return err
	}

	return nil
}

func getProperty(id uint) (*model.Property, error) {
	existingProperty := new(model.Property)
	result := db.DB.First(existingProperty, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", existingProperty)
	return existingProperty, nil
}

func mapToProtoProperty(property model.Property) *proto.Property {
	return &proto.Property{
		Id:          uint32(property.ID),
		Name:        property.Name,
		Description: property.Description,
		OwnerName:   property.OwnerName,
		Status:      string(property.Status),
		CreatedAt:   timestamppb.New(property.CreatedAt),
		UpdatedAt:   timestamppb.New(property.UpdatedAt),
	}
}
