package service

import (
	"context"
	"errors"

	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/producerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/secondaryport/repositoryport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

var (
	ErrErrorInjection = errors.New("error injected")
)

type HotelService struct {
	ErrorFlag       bool
	hotelProducer   producerport.HotelProducer
	hotelRepository repositoryport.HotelRepository
}

func NewHotelService(
	hotelProducer producerport.HotelProducer,
	hotelRepository repositoryport.HotelRepository,
) *HotelService {
	return &HotelService{
		ErrorFlag:       false,
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
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	if s.ErrorFlag {
		logger.Errorw("Error injection", "err", ErrErrorInjection, "booking", booking)
		span.RecordError(ErrErrorInjection)
		span.SetStatus(o11y.GetStatus(ErrErrorInjection))
		return ErrErrorInjection
	}

	if err := s.hotelProducer.PublishHotelBooked(ctx, &booking); err != nil {
		logger.Errorw("Failed to publish HotelBooked", "booking", booking, "err", err)
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}

	return nil
}

func (s *HotelService) InjectError(ctx context.Context, flag bool) error {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "HotelService",
		"method", "ToggleErrorInjection",
	)
	logger.Debugw("Toggle error")

	s.ErrorFlag = flag

	return nil
}

func (s *HotelService) GetErrorFlag(ctx context.Context) bool {
	logger := util.GetLogger().WithContext(ctx).With(
		"service", "HotelService",
		"method", "GetErrorFlag",
	)
	logger.Debugw("Get error flag")

	return s.ErrorFlag
}
