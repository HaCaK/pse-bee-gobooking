package handler

import (
	"context"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

// creates and starts a PropertyExternalServer and returns a client that is connected to it and can be used for tests
func startPropertyExternalServer(ctx context.Context) (proto.PropertyExternalClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	proto.RegisterPropertyExternalServer(baseServer, new(PropertyHandler))
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("Error serving propertyExternalServer: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Error connecting to propertyExternalServer: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("Error closing propertyExternalServer listener: %v", err)
		}
		baseServer.Stop()
	}

	client := proto.NewPropertyExternalClient(conn)

	return client, closer
}

func getMockListPropertiesResp(propertyResp *proto.PropertyResp) *proto.ListPropertiesResp {
	list := new(proto.ListPropertiesResp)

	if propertyResp != nil {
		list.Properties = append(list.Properties, propertyResp)
	}

	return list
}

func getMockPropertyRespWithDefaultOwnerName() *proto.PropertyResp {
	return getMockPropertyResp("owner")
}

func getMockPropertyResp(ownerName string) *proto.PropertyResp {
	return &proto.PropertyResp{
		Id:          1,
		Description: "description",
		OwnerName:   ownerName,
		Status:      "FREE",
	}
}

func createPropertyInDB() {
	property := model.Property{
		Description: "description",
		OwnerName:   "owner",
		Status:      "FREE",
	}
	db.DB.Create(&property)
}

func deletePropertyInDB() {
	db.DB.Delete(new(model.Property), 1)
}
