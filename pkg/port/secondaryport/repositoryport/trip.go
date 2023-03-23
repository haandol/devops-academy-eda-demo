package repositoryport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
)

type TripRepository interface {
	Create(ctx context.Context, tripID string) (dto.Trip, error)
	Complete(ctx context.Context, evt *event.FlightBooked) error
	List(ctx context.Context) ([]dto.Trip, error)
}
