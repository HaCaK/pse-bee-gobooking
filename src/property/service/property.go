package service

import (
	"errors"
	"github.com/HaCaK/pse-bee-gobooking/src/property/db"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateProperty(property *model.Property) error {
	property.SetStatusFree()

	result := db.DB.Create(property)
	if result.Error != nil {
		return result.Error
	}
	log.Infof("Successfully stored new property with ID %v in database.", property.ID)
	log.Tracef("Stored: %v", property)
	return nil
}

func GetProperties() ([]model.Property, error) {
	var properties []model.Property
	result := db.DB.Find(&properties)
	if result.Error != nil {
		return nil, result.Error
	}
	log.Tracef("Retrieved: %v", properties)
	return properties, nil
}

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

func DeleteProperty(id uint) (*model.Property, error) {
	existingProperty, err := GetProperty(id)
	if existingProperty == nil || err != nil {
		return existingProperty, err
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

func BookProperty(existingProperty *model.Property, bookingId uint) error {
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

func FreeProperty(existingProperty *model.Property) error {
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
