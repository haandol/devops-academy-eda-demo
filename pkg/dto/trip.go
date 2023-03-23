package dto

type Trip struct {
	ID              string `json:"id" binding:"required" validate:"required"`
	UserID          string `json:"userId" binding:"required" validate:"required"`
	CarID           string `json:"carId" binding:"required" validate:"required"`
	HotelID         string `json:"hotelId" binding:"required" validate:"required"`
	FlightID        string `json:"flightId" binding:"required" validate:"required"`
	CarBookingID    string `json:"carBookingId"`
	HotelBookingID  string `json:"hotelBookingId"`
	FlightBookingID string `json:"flightBookingId"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt       string `json:"updatedAt"`
}
