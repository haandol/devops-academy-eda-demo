package producerport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type TripProducer interface {
	PublishTripCreated(ctx context.Context, t *dto.Trip) error
}
