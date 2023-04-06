package event

import "github.com/haandol/devops-academy-eda-demo/pkg/message"

type TripCreated struct {
	message.Message
	Body TripCreatedBody `json:"body" validate:"required"`
}

type TripCreatedBody struct {
	TripID string `json:"tripId" validate:"required"`
}
