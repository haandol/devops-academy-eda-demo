package service

import (
	"context"
	"fmt"

	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type CarService struct {
	carProducer   producerport.CarProducer
	carRepository repositoryport.CarRepository
}

func NewCarService(
	carProducer producerport.CarProducer,
	carRepository repositoryport.CarRepository,
) *CarService {
	return &CarService{
		carProducer:   carProducer,
		carRepository: carRepository,
	}
}

func (s *CarService) Book(ctx context.Context, evt *event.TripCreated) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "CarService",
		"method", "Book",
	)
	logger.Debugw("Book car", "event", evt)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	span.SetAttributes(
		o11y.AttrString("evt", fmt.Sprintf("%v", evt)),
	)

	booking, err := s.carRepository.Book(ctx, evt.Body.TripID)
	if err != nil {
		logger.Errorw("Failed to book car", "err", err)
		return err
	}
	span.SetAttributes(
		o11y.AttrString("booking", fmt.Sprintf("%v", booking)),
	)

	if err := s.carProducer.PublishCarBooked(ctx, &booking); err != nil {
		logger.Errorw("Failed to publish CarBooked", "booking", booking, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}
