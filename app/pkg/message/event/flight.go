package event

import "github.com/haandol/devops-academy-eda-demo/pkg/message"

type FlightBooked struct {
	message.Message
	Body FlightBookedBody `json:"body" validate:"required"`
}

type FlightBookedBody struct {
	TripID    string `json:"tripId" validate:"required"`
	BookingID string `json:"bookingId" validate:"required"`
	FlightID  string `json:"flightId" validate:"required"`
}
