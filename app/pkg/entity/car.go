package entity

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type CarBooking struct {
	PK        string `json:"PK" binding:"required" validate:"required"`
	SK        string `json:"SK" binding:"required" validate:"required"`
	ID        string `json:"id"`
	TripID    string `json:"tripId" binding:"required" validate:"required"`
	CarID     string `json:"carId" binding:"required" validate:"required"`
	Status    string `json:"status" binding:"required" validate:"required"`
	CreatedAt string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}

func (c *CarBooking) DTO() dto.CarBooking {
	return dto.CarBooking{
		ID:        c.ID,
		TripID:    c.TripID,
		CarID:     c.CarID,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
