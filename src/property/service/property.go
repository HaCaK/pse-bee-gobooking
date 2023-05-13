package service

import (
	"errors"
	"fmt"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateProperty creates the given property with initial status FREE
func CreateProperty(property *model.Property) error {
	property.SetStatusFree()

	result := db.DB.Create(property)
	if result.Error != nil {
		return result.Error
	}
	entry := log.WithField("ID", property.ID)
	entry.Info("Successfully stored new property in database.")
	entry.Tracef("Stored: %v", property)
	return nil
}

// GetProperties retrieves all existing properties
func GetProperties() ([]model.Property, error) {
	var properties []model.Property
	result := db.DB.Find(&properties)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", properties)
	return properties, nil
}

// GetProperty retrieves the property matching the given id
func GetProperty(id uint) (*model.Property, error) {
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

// UpdateProperty updates the property matching the given id
func UpdateProperty(id uint, property *model.Property) (*model.Property, error) {
	existingProperty, err := GetProperty(id)
	if existingProperty == nil || err != nil {
		return existingProperty, err
	}

	existingProperty.Name = property.Name
	existingProperty.Description = property.Description
	existingProperty.OwnerName = property.OwnerName
	existingProperty.Address = property.Address

	result := db.DB.Save(existingProperty)
	if result.Error != nil {
		return nil, result.Error
	}

	entry := log.WithField("ID", id)
	entry.Info("Successfully updated property.")
	entry.Tracef("Updated: %v", existingProperty)
	return existingProperty, nil
}

// DeleteProperty deletes the property matching the given id
// NOTE: Deletion is only possible if the property is free
func DeleteProperty(id uint) (*model.Property, error) {
	existingProperty, err := GetProperty(id)
	if existingProperty == nil || err != nil {
		return existingProperty, err
	}

	if existingProperty.IsStatusBooked() {
		return nil, &model.PropertyError{Message: "Property cannot be deleted, because it is booked. Please, cancel the booking first."}
	}

	result := db.DB.Delete(existingProperty)
	if result.Error != nil {
		return nil, result.Error
	}

	entry := log.WithField("ID", id)
	entry.Info("Successfully deleted property.")
	entry.Tracef("Deleted: %v", existingProperty)
	return existingProperty, nil
}

// BookProperty books the given property if it is not already booked
// This is checked to prevent double-booking the property
func BookProperty(existingProperty *model.Property, bookingId uint) error {
	if existingProperty.IsStatusBooked() {
		message := fmt.Sprintf("Sorry, property %s (ID: %d) is already booked", existingProperty.Name, existingProperty.ID)
		return &model.PropertyError{Message: message}
	}

	existingProperty.SetStatusBooked()
	existingProperty.BookingId = bookingId

	result := db.DB.Save(existingProperty)
	if result.Error != nil {
		return result.Error
	}

	entry := log.WithField("ID", existingProperty.ID)
	entry.Info("Successfully booked property.")
	entry.Tracef("Updated: %v", existingProperty)
	return nil
}

// FreeProperty frees the given property if the given requestedBookingId matches the stored bookingId
// This is checked to prevent someone from cancelling another person's booking
func FreeProperty(existingProperty *model.Property, requestedBookingId uint) error {
	if existingProperty.BookingId != requestedBookingId {
		message := fmt.Sprintf("Whoops! It seems as if the property %s (ID: %d) is already booked.", existingProperty.Name, existingProperty.ID)
		return &model.PropertyError{Message: message}
	}

	existingProperty.SetStatusFree()
	existingProperty.BookingId = 0

	result := db.DB.Save(existingProperty)
	if result.Error != nil {
		return result.Error
	}

	entry := log.WithField("ID", existingProperty.ID)
	entry.Info("Successfully freed property.")
	entry.Tracef("Updated: %v", existingProperty)
	return nil
}
