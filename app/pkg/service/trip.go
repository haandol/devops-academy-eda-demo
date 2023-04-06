package service

import (
	"context"
	"fmt"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/restport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type TripService struct {
	tripProducer    producerport.TripProducer
	tripRepository  repositoryport.TripRepository
	tripRestAdapter restport.TripRestAdapter
}

func NewTripService(
	tripProducer producerport.TripProducer,
	tripRepository repositoryport.TripRepository,
	tripRestAdapter restport.TripRestAdapter,
) *TripService {
	return &TripService{
		tripProducer:    tripProducer,
		tripRepository:  tripRepository,
		tripRestAdapter: tripRestAdapter,
	}
}

// create trip for the given user
func (s *TripService) Create(ctx context.Context, tripID string) (dto.Trip, error) {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "Create",
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Create")
	defer span.End()

	span.SetAttributes(
		o11y.AttrString("tripID", tripID),
	)

	trip, err := s.tripRepository.Create(ctx, tripID)
	if err != nil {
		logger.Errorw("failed to create trip", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return dto.Trip{}, err
	}
	span.SetAttributes(
		o11y.AttrString("trip", fmt.Sprintf("%v", trip)),
	)

	logger.Infow("publishing TripCreated", "trip", trip, "producer", s.tripProducer)
	if err := s.tripProducer.PublishTripCreated(ctx, &trip); err != nil {
		logger.Errorw("failed to publish TripCreated", "trip", trip, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return dto.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) Complete(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "Complete",
	)

	ctx, span := o11y.BeginSubSpan(ctx, "Complete")
	defer span.End()

	span.SetAttributes(
		o11y.AttrString("evt", fmt.Sprintf("%v", evt)),
	)

	if err := s.tripRepository.Complete(ctx, evt); err != nil {
		logger.Errorw("failed to update trip booking", "event", evt, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *TripService) List(ctx context.Context) ([]dto.Trip, error) {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "List",
	)

	trips, err := s.tripRepository.List(ctx)
	if err != nil {
		logger.Errorw("failed to create trip", "err", err)
		return []dto.Trip{}, err
	}

	return trips, nil
}

func (s *TripService) GetInjectionStatus(ctx context.Context) (bool, error) {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "GetInjectionStatus",
	)

	injectionStatus, err := s.tripRestAdapter.GetInjectionStatus(ctx)
	if err != nil {
		logger.Errorw("failed to get error injection status from hotel", "err", err)
		return false, err
	}

	return injectionStatus, nil
}

func (s *TripService) InjectError(ctx context.Context, flag bool) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "InjectError",
	)

	ctx, span := o11y.BeginSubSpan(ctx, "InjectError")
	defer span.End()

	span.SetAttributes(
		o11y.AttrBool("flag", flag),
	)

	if err := s.tripRestAdapter.InjectError(ctx, flag); err != nil {
		logger.Errorw("failed to inject error to hotel", "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}
