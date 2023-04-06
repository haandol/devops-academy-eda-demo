package event

import "github.com/haandol/devops-academy-eda-demo/pkg/message"

type HotelBooked struct {
	message.Message
	Body HotelBookedBody `json:"body" validate:"required"`
}

type HotelBookedBody struct {
	TripID    string `json:"tripId" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
	HotelID   string `json:"hotelId" validate:"required"`
}
