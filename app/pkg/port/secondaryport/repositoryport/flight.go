package repositoryport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type FlightRepository interface {
	Book(ctx context.Context, tripID string) (dto.FlightBooking, error)
}
