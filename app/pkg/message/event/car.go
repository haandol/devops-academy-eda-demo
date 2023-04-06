package event

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/message"
)

type CarBooked struct {
	message.Message
	Body CarBookedBody `json:"body" validate:"required"`
}

type CarBookedBody struct {
	TripID    string `json:"tripId" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
	CarID     string `json:"carId" validate:"required"`
}
