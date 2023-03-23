package repositoryport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type HotelRepository interface {
	Book(ctx context.Context, tripID string) (dto.HotelBooking, error)
}
