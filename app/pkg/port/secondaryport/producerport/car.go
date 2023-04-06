package producerport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type CarProducer interface {
	PublishCarBooked(ctx context.Context, d *dto.CarBooking) error
}
