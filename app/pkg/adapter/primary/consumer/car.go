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

type CarConsumer struct {
	*KafkaConsumer
	carService *service.CarService
}

func NewCarConsumer(
	kafkaConsumer *KafkaConsumer,
	carService *service.CarService,
) *CarConsumer {
	return &CarConsumer{
		KafkaConsumer: kafkaConsumer,
		carService:    carService,
	}
}

func (c *CarConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "CarConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err)
	}
}

func (c *CarConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "CarConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal message", "err", err)
	}

	logger.Infow("Received event", "event", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "CarConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "TripCreated":
		evt := &event.TripCreated{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Errorw("Failed to unmarshal event", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.carService.Book(ctx, evt)
	default:
		logger.Errorw("unknown event", "message", msg)
		err := errors.New("unknown event")
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
