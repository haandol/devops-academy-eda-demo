package service

import (
	"context"

	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type HotelService struct {
	hotelProducer   producerport.HotelProducer
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	hotelProducer producerport.HotelProducer,
	hotelRepository repositoryport.HotelRepository,
) *HotelService {
	return &HotelService{
		hotelProducer:   hotelProducer,
		hotelRepository: hotelRepository,
	}
}

func (s *HotelService) Book(ctx context.Context, evt *event.CarBooked) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "HotelService",
		"method", "Book",
	)
	logger.Debugw("Book hotel", "event", evt)

	ctx, span := o11y.BeginSubSpan(ctx, "Book")
	defer span.End()

	booking, err := s.hotelRepository.Book(ctx, evt.Body.TripID)
	if err != nil {
		logger.Errorw("Failed to book hotel", "err", err)
		return err
	}

	if err := s.hotelProducer.PublishHotelBooked(ctx, &booking); err != nil {
		logger.Errorw("Failed to publish HotelBooked", "booking", booking, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}
