package handler

import (
	"context"
	"errors"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/proto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

type PropertyTestSuite struct {
	suite.Suite
	cleanUpDB func()
}

// beforeAll
func (suite *PropertyTestSuite) SetupSuite() {
	log.Info(">>> From SetupSuite")
	suite.cleanUpDB = db.SetupTestDB(suite.T())
}

// afterAll
func (suite *PropertyTestSuite) TearDownSuite() {
	log.Info(">>> From TearDownSuite")
	suite.cleanUpDB()
}

func (suite *PropertyTestSuite) TestPropertyHandler_GetProperties() {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *proto.ListPropertiesResp
		err error
	}

	tests := map[string]struct {
		in           *emptypb.Empty
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoProperty_WhenGetProperties_ThenReturnEmpty": {
			in:           new(emptypb.Empty),
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: getMockListPropertiesResp(nil),
				err: nil,
			},
		},
		"GivenOneProperty_WhenGetProperties_ThenReturnProperty": {
			in: new(emptypb.Empty),
			setupFunc: func() {
				createPropertyInDB()
			},
			tearDownFunc: func() {
				deletePropertyInDB()
			},
			expected: expectation{
				out: getMockListPropertiesResp(getMockPropertyRespWithDefaultOwnerName()),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := client.GetProperties(ctx, testData.in)
		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			len(out.Properties) != len(testData.expected.out.Properties) ||
			(len(out.Properties) > 0 && out.Properties[0].OwnerName != testData.expected.out.Properties[0].OwnerName) {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *PropertyTestSuite) TestPropertyHandler_GetProperty() {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *proto.PropertyResp
		err error
	}

	tests := map[string]struct {
		in           *proto.PropertyIdReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoPropertyAndEmptyPropertyIdReq_WhenGetProperty_ThenReturnNotFound": {
			in:           &proto.PropertyIdReq{},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenNoProperty_WhenGetProperty_ThenReturnNotFound": {
			in:           &proto.PropertyIdReq{Id: 1},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenOneProperty_WhenGetProperty_ThenReturnProperty": {
			in: &proto.PropertyIdReq{Id: 1},
			setupFunc: func() {
				createPropertyInDB()
			},
			tearDownFunc: func() {
				deletePropertyInDB()
			},
			expected: expectation{
				out: getMockPropertyRespWithDefaultOwnerName(),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := client.GetProperty(ctx, testData.in)
		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			out.OwnerName != testData.expected.out.OwnerName {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *PropertyTestSuite) TestPropertyHandler_UpdateProperty() {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *proto.PropertyResp
		err error
	}

	tests := map[string]struct {
		in           *proto.UpdatePropertyReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoPropertyAndEmptyUpdatePropertyReq_WhenUpdateProperty_ThenReturnNotFound": {
			in:           &proto.UpdatePropertyReq{},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenNoProperty_WhenUpdateProperty_ThenReturnNotFound": {
			in:           &proto.UpdatePropertyReq{Id: 1, OwnerName: "other"},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenOneProperty_WhenUpdateProperty_ThenReturnUpdatedProperty": {
			in: &proto.UpdatePropertyReq{Id: 1, OwnerName: "other"},
			setupFunc: func() {
				createPropertyInDB()
			},
			tearDownFunc: func() {
				deletePropertyInDB()
			},
			expected: expectation{
				out: getMockPropertyResp("other"),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := client.UpdateProperty(ctx, testData.in)

		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			testData.expected.out.OwnerName != out.OwnerName {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *PropertyTestSuite) TestPropertyHandler_CreateProperty() {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *proto.PropertyResp
		err error
	}

	tests := map[string]struct {
		in           *proto.CreatePropertyReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoProperty_WhenCreateProperty_ThenReturnCreatedProperty": {
			in:        &proto.CreatePropertyReq{OwnerName: "owner"},
			setupFunc: nil,
			tearDownFunc: func() {
				deletePropertyInDB()
			},
			expected: expectation{
				out: getMockPropertyResp("owner"),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := client.CreateProperty(ctx, testData.in)

		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			testData.expected.out.OwnerName != out.OwnerName {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *PropertyTestSuite) TestPropertyHandler_DeleteProperty() {
	ctx := context.Background()
	client, closer := server(ctx)
	defer closer()

	type expectation struct {
		out *emptypb.Empty
		err error
	}

	tests := map[string]struct {
		in           *proto.PropertyIdReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoPropertyAndEmptyPropertyIdReq_WhenDeleteProperty_ThenReturnNotFound": {
			in:           &proto.PropertyIdReq{},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenNoProperty_WhenDeleteProperty_ThenReturnNotFound": {
			in:           &proto.PropertyIdReq{Id: 1},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Property not found"),
			},
		},
		"GivenOneProperty_WhenDeleteProperty_ThenDeleteProperty": {
			in: &proto.PropertyIdReq{Id: 1},
			setupFunc: func() {
				createPropertyInDB()
			},
			tearDownFunc: func() {
				deletePropertyInDB()
			},
			expected: expectation{
				out: new(emptypb.Empty),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := client.DeleteProperty(ctx, testData.in)
		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func TestPropertyTestSuite(t *testing.T) {
	suite.Run(t, new(PropertyTestSuite))
}
