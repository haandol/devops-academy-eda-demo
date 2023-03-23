package producer

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/devops-academy-eda-demo/pkg/connector/producer"
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/message"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/haandol/devops-academy-eda-demo/pkg/util/o11y"
)

type FlightProducer struct {
	*producer.KafkaProducer
}

func NewFlightProducer(kafkaProducer *producer.KafkaProducer) *FlightProducer {
	return &FlightProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *FlightProducer) PublishFlightBooked(ctx context.Context, d *dto.FlightBooking) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	evt := &event.FlightBooked{
		Message: message.Message{
			Name:          reflect.ValueOf(event.FlightBooked{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.FlightBookedBody{
			TripID:    d.TripID,
			BookingID: d.ID,
			FlightID:  d.FlightID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "trip-service", d.TripID, v); err != nil {
		return err
	}

	return nil
}
