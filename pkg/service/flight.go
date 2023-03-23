package service

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type FlightService struct {
	flightProducer   producerport.FlightProducer
	flightRepository repositoryport.FlightRepository
}

func NewFlightService(
	flightProducer producerport.FlightProducer,
	flightRepository repositoryport.FlightRepository,
) *FlightService {
	return &FlightService{
		flightProducer:   flightProducer,
		flightRepository: flightRepository,
	}
}

func (s *FlightService) Book(ctx context.Context, evt *event.HotelBooked) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "FlightService",
		"method", "Book",
		"event", evt,
	)
	logger.Infow("Book flight", "event", evt)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	booking, err := s.flightRepository.Book(ctx, evt.Body.TripID)
	if err != nil {
		logger.Errorw("failed to book car", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	if err := s.flightProducer.PublishFlightBooked(ctx, &booking); err != nil {
		logger.Errorw("failed to publish flight booked", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}
