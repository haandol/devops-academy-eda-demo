package entity

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type HotelBooking struct {
	PK        string `json:"PK" binding:"required" validate:"required"`
	SK        string `json:"SK" binding:"required" validate:"required"`
	ID        string `json:"id"`
	TripID    string `json:"tripId" binding:"required" validate:"required"`
	HotelID   string `json:"hotelId" binding:"required" validate:"required"`
	Status    string `json:"status" binding:"required" validate:"required"`
	CreatedAt string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}

func (h *HotelBooking) DTO() dto.HotelBooking {
	return dto.HotelBooking{
		ID:        h.ID,
		TripID:    h.TripID,
		HotelID:   h.HotelID,
		Status:    h.Status,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}
