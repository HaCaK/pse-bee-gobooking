package model

import "gorm.io/gorm"

type Status string

const (
	PENDING   Status = "PENDING"
	CONFIRMED Status = "CONFIRMED"
)

type Booking struct {
	gorm.Model
	Comment      string `gorm:"notNull;size:100"`
	CustomerName string `gorm:"notNull;size:60"`
	Status       `gorm:"notNull;type:ENUM('PENDING', 'CONFIRMED')"`
	PropertyId   uint `gorm:"notNull"`
}

func (booking *Booking) SetStatusPending() {
	booking.Status = PENDING
}

func (booking *Booking) SetStatusConfirmed() {
	booking.Status = CONFIRMED
}
