package model

import "gorm.io/gorm"

type Status string

const (
	FREE   Status = "FREE"
	BOOKED Status = "BOOKED"
)

type Property struct {
	gorm.Model
	Name        string `gorm:"notNull;size:60"`
	Description string `gorm:"notNull;size:100"`
	OwnerName   string `gorm:"notNull;size:60"`
	Address     string `gorm:"notNull;size:100"`
	BookingId   uint
	Status      `gorm:"notNull;type:ENUM('FREE', 'BOOKED')"`
}

func (property *Property) SetStatusFree() {
	property.Status = FREE
}

func (property *Property) SetStatusBooked() {
	property.Status = BOOKED
}

func (property *Property) IsStatusBooked() bool {
	return property.Status == BOOKED
}
