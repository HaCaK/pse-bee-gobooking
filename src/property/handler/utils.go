package handler

import (
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapToProtoPropertyResp(property *model.Property) *proto.PropertyResp {
	return &proto.PropertyResp{
		Id:          uint32(property.ID),
		Name:        property.Name,
		Description: property.Description,
		OwnerName:   property.OwnerName,
		Status:      string(property.Status),
		CreatedAt:   timestamppb.New(property.CreatedAt),
		UpdatedAt:   timestamppb.New(property.UpdatedAt),
	}
}
