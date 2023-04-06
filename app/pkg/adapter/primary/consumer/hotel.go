package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/haandol/devops-academy-eda-demo/pkg/message"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/consumerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/service"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type HotelConsumer struct {
	*KafkaConsumer
	hotelService *service.HotelService
}

func NewHotelConsumer(
	kafkaConsumer *KafkaConsumer,
	hotelService *service.HotelService,
) *HotelConsumer {
	return &HotelConsumer{
		KafkaConsumer: kafkaConsumer,
		hotelService:  hotelService,
	}
}

func (c *HotelConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "HotelConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err)
	}
}

func (c *HotelConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "HotelConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal event", "err", err)
	}

	logger.Infow("Received event", "event", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "HotelConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "CarBooked":
		evt := &event.CarBooked{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Errorw("Failed to unmarshal event", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.hotelService.Book(ctx, evt)
	default:
		logger.Errorw("unknown event", "message", msg)
		err := errors.New("unknown event")
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
