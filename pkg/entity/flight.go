package entity

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type FlightBooking struct {
	PK        string `json:"PK" binding:"required" validate:"required"`
	SK        string `json:"SK" binding:"required" validate:"required"`
	ID        string `json:"id"`
	TripID    string `json:"tripId" binding:"required" validate:"required"`
	FlightID  string `json:"flightId" binding:"required" validate:"required"`
	Status    string `json:"status" binding:"required" validate:"required"`
	CreatedAt string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}

func (f *FlightBooking) DTO() dto.FlightBooking {
	return dto.FlightBooking{
		ID:        f.ID,
		TripID:    f.TripID,
		FlightID:  f.FlightID,
		Status:    f.Status,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}
