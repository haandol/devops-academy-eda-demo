package dto

type FlightBooking struct {
	ID        string `json:"id"`
	TripID    string `json:"tripId" binding:"required" validate:"required"`
	FlightID  string `json:"flightId" binding:"required" validate:"required"`
	Status    string `json:"status" binding:"required" validate:"required"`
	CreatedAt string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}
