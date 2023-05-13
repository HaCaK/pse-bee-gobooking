package handler

import (
	"context"
	"errors"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/db"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/handler/integration_test"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

const propertyInternalServerPort = "9111"

type BookingTestSuite struct {
	suite.Suite
	ctx                        context.Context
	client                     proto.BookingExternalClient
	closeBookingExternalServer func()
	mockPropertyInternalServer *integration_test.MockPropertyInternalServer
	cleanUpDB                  func()
}

// beforeAll
func (suite *BookingTestSuite) SetupSuite() {
	log.Info(">>> From SetupSuite")
	suite.ctx = context.Background()
	suite.client, suite.closeBookingExternalServer = startBookingExternalServer(suite.ctx)
	suite.mockPropertyInternalServer = new(integration_test.MockPropertyInternalServer)
}

// beforeEach
func (suite *BookingTestSuite) SetupTest() {
	log.Info("--- From SetupTest: Setting up fresh DB")
	suite.cleanUpDB = db.SetupTestDB(suite.T())
}

// afterAll
func (suite *BookingTestSuite) TearDownSuite() {
	log.Info(">>> From TearDownSuite")
	suite.closeBookingExternalServer()
}

// afterEach
func (suite *BookingTestSuite) TearDownTest() {
	log.Info("--- From TearDownTest: Cleaning up DB")
	suite.cleanUpDB()
}

func (suite *BookingTestSuite) TestBookingHandler_GetBookings() {
	type expectation struct {
		out *proto.ListBookingsResp
		err error
	}

	tests := map[string]struct {
		in           *emptypb.Empty
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoBooking_WhenGetBookings_ThenReturnEmpty": {
			in:           new(emptypb.Empty),
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: getMockListBookingsResp(nil),
				err: nil,
			},
		},
		"GivenOneBooking_WhenGetBookings_ThenReturnBooking": {
			in: new(emptypb.Empty),
			setupFunc: func() {
				createBookingInDB()
			},
			tearDownFunc: func() {
				deleteBookingInDB()
			},
			expected: expectation{
				out: getMockListBookingsResp(getMockBookingRespWithDefaultCustomerName()),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := suite.client.GetBookings(suite.ctx, testData.in)
		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			len(out.Bookings) != len(testData.expected.out.Bookings) ||
			(len(out.Bookings) > 0 && out.Bookings[0].CustomerName != testData.expected.out.Bookings[0].CustomerName) {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *BookingTestSuite) TestBookingHandler_GetBooking() {
	type expectation struct {
		out *proto.BookingResp
		err error
	}

	tests := map[string]struct {
		in           *proto.BookingIdReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoBookingAndEmptyBookingIdReq_WhenGetBooking_ThenReturnNotFound": {
			in:           &proto.BookingIdReq{},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Booking not found"),
			},
		},
		"GivenNoBooking_WhenGetBooking_ThenReturnNotFound": {
			in:           &proto.BookingIdReq{Id: 1},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Booking not found"),
			},
		},
		"GivenOneBooking_WhenGetBooking_ThenReturnBooking": {
			in: &proto.BookingIdReq{Id: 1},
			setupFunc: func() {
				createBookingInDB()
			},
			tearDownFunc: func() {
				deleteBookingInDB()
			},
			expected: expectation{
				out: getMockBookingRespWithDefaultCustomerName(),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := suite.client.GetBooking(suite.ctx, testData.in)
		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			out.CustomerName != testData.expected.out.CustomerName {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *BookingTestSuite) TestBookingHandler_UpdateBooking() {
	type expectation struct {
		out *proto.BookingResp
		err error
	}

	tests := map[string]struct {
		in           *proto.UpdateBookingReq
		setupFunc    func()
		tearDownFunc func()
		expected     expectation
	}{
		"GivenNoBookingAndEmptyUpdateBookingReq_WhenUpdateBooking_ThenReturnNotFound": {
			in:           &proto.UpdateBookingReq{},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Booking not found"),
			},
		},
		"GivenNoBooking_WhenUpdateBooking_ThenReturnNotFound": {
			in:           &proto.UpdateBookingReq{Id: 1, CustomerName: "other"},
			setupFunc:    nil,
			tearDownFunc: nil,
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = NotFound desc = Booking not found"),
			},
		},
		"GivenOneBooking_WhenUpdateBooking_ThenReturnUpdatedBooking": {
			in: &proto.UpdateBookingReq{Id: 1, CustomerName: "other"},
			setupFunc: func() {
				createBookingInDB()
			},
			tearDownFunc: func() {
				deleteBookingInDB()
			},
			expected: expectation{
				out: getMockBookingResp("other"),
				err: nil,
			},
		},
	}

	for scenario, testData := range tests {
		log.Infof("Scenario: %s", scenario)

		if testData.setupFunc != nil {
			testData.setupFunc()
		}

		out, err := suite.client.UpdateBooking(suite.ctx, testData.in)

		if err != nil {
			if testData.expected.err.Error() != err.Error() {
				suite.T().Errorf("Err:\n Expected: %v\n Actual: %v", testData.expected.err, err)
			}
		} else if out == nil ||
			testData.expected.out.CustomerName != out.CustomerName {
			suite.T().Errorf("Out:\n Expected: %v\n Actual: %v", testData.expected.out, out)
		}

		if testData.tearDownFunc != nil {
			testData.tearDownFunc()
		}
	}
}

func (suite *BookingTestSuite) TestBookingHandler_CreateBooking() {
	cancel := suite.mockPropertyInternalServer.Start(propertyInternalServerPort)
	defer cancel()

	// given
	defer deleteBookingInDB()
	in := &proto.CreateBookingReq{
		Comment:      "test",
		CustomerName: "cust",
		PropertyId:   1,
	}

	// when
	out, err := suite.client.CreateBooking(suite.ctx, in)

	// then
	if err != nil {
		suite.T().Errorf("Unexpected err: %v", err)
	} else if out == nil || out.CustomerName != "cust" || out.Status != "CONFIRMED" {
		suite.T().Errorf("Unexpected: %v", out)
	}
}

func (suite *BookingTestSuite) TestBookingHandler_DeleteBooking() {
	cancel := suite.mockPropertyInternalServer.Start(propertyInternalServerPort)
	defer cancel()

	// given
	createBookingInDB()
	defer deleteBookingInDB()

	in := &proto.BookingIdReq{
		Id: 1,
	}

	// when
	_, err := suite.client.DeleteBooking(suite.ctx, in)

	// then
	if err != nil {
		suite.T().Errorf("Unexpected err: %v", err)
	}
}

func TestBookingTestSuite(t *testing.T) {
	suite.Run(t, new(BookingTestSuite))
}
