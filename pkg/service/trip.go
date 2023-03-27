package service

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/restport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
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

	trip, err := s.tripRepository.Create(ctx, tripID)
	if err != nil {
		logger.Errorw("failed to create trip", "err", err)
		return dto.Trip{}, err
	}

	logger.Infow("publishing TripCreated", "trip", trip, "producer", s.tripProducer)
	if err := s.tripProducer.PublishTripCreated(ctx, &trip); err != nil {
		logger.Errorw("failed to publish TripCreated", "trip", trip, "err", err)
		return dto.Trip{}, err
	}

	return trip, nil
}

func (s *TripService) Complete(ctx context.Context, evt *event.FlightBooked) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "TripService",
		"method", "Complete",
	)

	if err := s.tripRepository.Complete(ctx, evt); err != nil {
		logger.Errorw("failed to update trip booking", "event", evt, "err", err)
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

	if err := s.tripRestAdapter.InjectError(ctx, flag); err != nil {
		logger.Errorw("failed to inject error to hotel", "err", err)
		return err
	}

	return nil
}
