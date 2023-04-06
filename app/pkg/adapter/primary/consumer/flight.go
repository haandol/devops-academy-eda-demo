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

type FlightConsumer struct {
	*KafkaConsumer
	flightService *service.FlightService
}

func NewFlightConsumer(
	kafkaConsumer *KafkaConsumer,
	flightService *service.FlightService,
) *FlightConsumer {
	return &FlightConsumer{
		KafkaConsumer: kafkaConsumer,
		flightService: flightService,
	}
}

func (c *FlightConsumer) Init() {
	logger := util.GetLogger().With(
		"module", "FlightConsumer",
		"func", "Init",
	)

	if err := c.RegisterHandler(c.Handle); err != nil {
		logger.Panicw("Failed to register handler", "err", err)
	}
}

func (c *FlightConsumer) Handle(ctx context.Context, r *consumerport.Message) error {
	logger := util.GetLogger().With(
		"module", "FlightConsumer",
		"func", "Handle",
	)

	msg := &message.Message{}
	if err := json.Unmarshal(r.Value, msg); err != nil {
		logger.Errorw("Failed to unmarshal command", "err", err)
	}

	logger.Infow("Received command", "command", msg)
	ctx, span := o11y.BeginSpanWithTraceID(ctx, msg.CorrelationID, msg.ParentID, "FlightConsumer")
	defer span.End()
	span.SetAttributes(
		o11y.AttrString("msg", fmt.Sprintf("%v", msg)),
	)

	switch msg.Name {
	case "HotelBooked":
		evt := &event.HotelBooked{}
		if err := json.Unmarshal(r.Value, evt); err != nil {
			logger.Errorw("Failed to unmarshal command", "err", err)
			span.RecordError(err)
			span.SetStatus(o11y.GetStatus(err))
			return err
		}
		return c.flightService.Book(ctx, evt)
	default:
		logger.Errorw("unknown command", "message", msg)
		err := errors.New("unknown command")
		span.RecordError(err)
		span.SetStatus(o11y.GetStatus(err))
		return err
	}
}
