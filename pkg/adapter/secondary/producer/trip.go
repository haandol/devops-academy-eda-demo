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

type TripProducer struct {
	*producer.KafkaProducer
}

func NewTripProducer(kafkaProducer *producer.KafkaProducer) *TripProducer {
	return &TripProducer{
		KafkaProducer: kafkaProducer,
	}
}

func (p *TripProducer) PublishTripCreated(ctx context.Context, trip *dto.Trip) error {
	traceID, spanID := o11y.GetTraceSpanID(ctx)
	evt := &event.TripCreated{
		Message: message.Message{
			Name:          reflect.ValueOf(event.TripCreated{}).Type().Name(),
			Version:       "1.0.0",
			ID:            uuid.NewString(),
			CorrelationID: traceID,
			ParentID:      spanID,
			CreatedAt:     time.Now().Format(time.RFC3339),
		},
		Body: event.TripCreatedBody{
			TripID: trip.ID,
		},
	}
	if err := util.ValidateStruct(evt); err != nil {
		return err
	}
	v, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	if err := p.Produce(ctx, "car-service", trip.ID, v); err != nil {
		return err
	}

	return nil
}
