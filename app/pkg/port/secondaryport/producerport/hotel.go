package producerport

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type HotelProducer interface {
	PublishHotelBooked(ctx context.Context, d *dto.HotelBooking) error
}
