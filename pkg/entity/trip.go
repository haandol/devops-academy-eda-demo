package entity

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type Trip struct {
	PK              string `json:"PK" binding:"required" validate:"required"`
	SK              string `json:"SK" binding:"required" validate:"required"`
	GS1PK           string `json:"GS1PK" binding:"required" validate:"required"`
	GS1SK           string `json:"GS1SK" binding:"required" validate:"required"`
	ID              string `json:"id" binding:"required" validate:"required"`
	CarID           string `json:"carId"`
	HotelID         string `json:"hotelId"`
	FlightID        string `json:"flightId"`
	CarBookingID    string `json:"carBookingId"`
	HotelBookingID  string `json:"hotelBookingId"`
	FlightBookingID string `json:"flightBookingId"`
	Status          string `json:"status" binding:"required" validate:"required"`
	CreatedAt       string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt       string `json:"updatedAt"`
}

func (t *Trip) DTO() dto.Trip {
	return dto.Trip{
		ID:              t.ID,
		CarID:           t.CarID,
		HotelID:         t.HotelID,
		FlightID:        t.FlightID,
		CarBookingID:    t.CarBookingID,
		HotelBookingID:  t.HotelBookingID,
		FlightBookingID: t.FlightBookingID,
		Status:          t.Status,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}
