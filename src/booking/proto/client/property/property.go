package client

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

var (
	propertyTarget = os.Getenv("PROPERTY_CONNECT")
)

func GetPropertyConnection(ctx context.Context) (*grpc.ClientConn, error) {
	var err error
	log.WithFields(log.Fields{
		"target": propertyTarget,
	}).Infoln("Connecting to property service")
	var conn *grpc.ClientConn
	conn, err = grpc.DialContext(ctx, propertyTarget, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, err
}
