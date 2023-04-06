package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/haandol/devops-academy-eda-demo/pkg/config"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant/status"
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/entity"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type FlightRepository struct {
	TableName string
	Client    *dynamodb.Client
}

func NewFlightRepository(
	cfg *config.Config,
	client *dynamodb.Client,
) *FlightRepository {
	return &FlightRepository{
		TableName: cfg.Database.TableName,
		Client:    client,
	}
}

func (r *FlightRepository) Book(ctx context.Context, tripID string) (dto.FlightBooking, error) {
	if booking, err := r.GetByTripID(ctx, tripID); err != nil {
		return dto.FlightBooking{}, err
	} else if booking.Status == status.Booked {
		return booking, nil
	}

	bookingID := uuid.NewString()
	req := &entity.FlightBooking{
		PK:        fmt.Sprintf("TRIP#%s", tripID),
		SK:        fmt.Sprintf("BOOKING#FLIGHT#%s", bookingID),
		ID:        bookingID,
		TripID:    tripID,
		FlightID:  uuid.NewString(),
		Status:    status.Booked,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := util.ValidateStruct(req); err != nil {
		return dto.FlightBooking{}, err
	}

	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		return dto.FlightBooking{}, err
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return dto.FlightBooking{}, err
	}

	return req.DTO(), nil
}

func (r *FlightRepository) GetByTripID(ctx context.Context, tripID string) (dto.FlightBooking, error) {
	pk, err := attributevalue.Marshal(fmt.Sprintf("TRIP#%s", tripID))
	if err != nil {
		return dto.FlightBooking{}, err
	}
	sk, err := attributevalue.Marshal(fmt.Sprintf("BOOKING#FLIGHT#%s", tripID))
	if err != nil {
		return dto.FlightBooking{}, err
	}

	res, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key:       map[string]types.AttributeValue{"PK": pk, "SK": sk},
	})
	if err != nil {
		return dto.FlightBooking{}, err
	}

	var row = &entity.FlightBooking{}
	err = attributevalue.UnmarshalMap(res.Item, row)
	if err != nil {
		return dto.FlightBooking{}, err
	}

	return row.DTO(), nil
}
