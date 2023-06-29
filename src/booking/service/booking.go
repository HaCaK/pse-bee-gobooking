package service

import (
	"context"
	"errors"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/db"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/model"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto"
	"github.com/HaCaK/pse-bee-gobooking/src/booking/proto/client/property"
	"google.golang.org/grpc"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateBooking creates the given booking
// and tries to confirm the booking at the property service
func CreateBooking(booking *model.Booking) error {
	booking.SetStatusPending()

	result := db.DB.Create(booking)
	if result.Error != nil {
		return result.Error
	}

	entry := log.WithField("ID", booking.ID)
	entry.Info("Successfully stored new booking in database.")
	entry.Tracef("Stored: %v", booking)

	err := confirmBooking(booking)
	if err != nil {
		return err
	}

	entry.Info("Successfully confirmed booking.")
	return nil
}

// GetBookings retrieves all existing bookings
func GetBookings() ([]model.Booking, error) {
	var bookings []model.Booking
	result := db.DB.Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", bookings)
	return bookings, nil
}

// GetBooking retrieves the booking matching the given id
func GetBooking(id uint) (*model.Booking, error) {
	booking := new(model.Booking)
	result := db.DB.First(booking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", booking)
	return booking, nil
}

// UpdateBooking updates the booking matching the given id
func UpdateBooking(id uint, booking *model.Booking) (*model.Booking, error) {
	existingBooking, err := GetBooking(id)
	if existingBooking == nil || err != nil {
		return existingBooking, err
	}
	existingBooking.CustomerName = booking.CustomerName
	existingBooking.Comment = booking.Comment

	result := db.DB.Save(existingBooking)
	if result.Error != nil {
		return nil, result.Error
	}

	entry := log.WithField("ID", id)
	entry.Info("Successfully updated booking.")
	entry.Tracef("Updated: %v", existingBooking)
	return existingBooking, nil
}

// DeleteBooking deletes the booking matching the given id
// and cancels it at the property service
func DeleteBooking(id uint) (*model.Booking, error) {
	booking, err := deleteWithoutCancellation(id)
	if booking == nil || err != nil {
		return booking, err
	}

	err = cancelBooking(booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

// deleteWithoutCancellation deletes the booking matching the given id
func deleteWithoutCancellation(id uint) (*model.Booking, error) {
	booking, err := GetBooking(id)
	if booking == nil || err != nil {
		return booking, err
	}
	result := db.DB.Delete(booking)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", id)
	entry.Info("Successfully deleted booking.")
	entry.Tracef("Deleted: %v", booking)

	return booking, nil
}

// confirmBooking connects to the property service via gRPC and confirms the given booking
// NOTE: Deletes the booking if the confirmation fails
func confirmBooking(booking *model.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := client.GetPropertyConnection(ctx)
	if err != nil {
		log.Errorf("Error connecting to property service: %v", err)
		return err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("Error closing connection: %s", err)
		}
	}(conn)

	propertyClient := proto.NewPropertyInternalClient(conn)
	_, err = propertyClient.ConfirmBooking(ctx, &proto.BookingReq{
		BookingId:  uint32(booking.ID),
		PropertyId: uint32(booking.PropertyId),
	})
	if err != nil {
		log.Errorf("Error calling property service: %v", err)
		entry := log.WithField("bookingId", booking.ID)
		entry.Info("Trying to delete booking to make state consistent")
		// does not require cancellation because booking was never confirmed
		_, deleteErr := deleteWithoutCancellation(booking.ID)

		if deleteErr != nil {
			return errors.Join(err, deleteErr)
		}
		return err
	}

	booking.SetStatusConfirmed()

	result := db.DB.Save(booking)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// cancelBooking connects to the property service via gRPC and cancels the given booking
func cancelBooking(booking *model.Booking) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, err := client.GetPropertyConnection(ctx)
	if err != nil {
		log.Errorf("Error connecting to property service: %v", err)
		return err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("Error closing connection: %s", err)
		}
	}(conn)

	propertyClient := proto.NewPropertyInternalClient(conn)
	_, err = propertyClient.CancelBooking(ctx, &proto.BookingReq{
		BookingId:  uint32(booking.ID),
		PropertyId: uint32(booking.PropertyId),
	})
	if err != nil {
		log.Errorf("Error calling property service: %v", err)
		return err
	}

	return nil
}
