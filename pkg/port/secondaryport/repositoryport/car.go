package repositoryport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type CarRepository interface {
	Book(ctx context.Context, tripID string) (dto.CarBooking, error)
}
