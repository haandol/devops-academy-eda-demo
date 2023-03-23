package producerport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type FlightProducer interface {
	PublishFlightBooked(ctx context.Context, d *dto.FlightBooking) error
}
