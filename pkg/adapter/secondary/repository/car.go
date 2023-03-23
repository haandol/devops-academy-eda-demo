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

type CarRepository struct {
	TableName string
	Client    *dynamodb.Client
}

func NewCarRepository(
	cfg *config.Config,
	client *dynamodb.Client,
) *CarRepository {
	return &CarRepository{
		TableName: cfg.Database.TableName,
		Client:    client,
	}
}

func (r *CarRepository) Book(ctx context.Context, tripID string) (dto.CarBooking, error) {
	if booking, err := r.GetByTripID(ctx, tripID); err != nil {
		return dto.CarBooking{}, err
	} else if booking.Status == status.Booked {
		return booking, nil
	}

	bookingID := uuid.NewString()
	req := &entity.CarBooking{
		PK:        fmt.Sprintf("TRIP#%s", tripID),
		SK:        fmt.Sprintf("BOOKING#CAR#%s", bookingID),
		ID:        bookingID,
		TripID:    tripID,
		CarID:     uuid.NewString(),
		Status:    status.Booked,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := util.ValidateStruct(req); err != nil {
		return dto.CarBooking{}, err
	}

	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		return dto.CarBooking{}, err
	}

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return dto.CarBooking{}, err
	}

	return req.DTO(), nil
}

func (r *CarRepository) GetByTripID(ctx context.Context, tripID string) (dto.CarBooking, error) {
	pk, err := attributevalue.Marshal(fmt.Sprintf("TRIP#%s", tripID))
	if err != nil {
		return dto.CarBooking{}, err
	}
	sk, err := attributevalue.Marshal("begins_with(SK, BOOKING#CAR#)")
	if err != nil {
		return dto.CarBooking{}, err
	}

	res, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key:       map[string]types.AttributeValue{"PK": pk, "SK": sk},
	})
	if err != nil {
		return dto.CarBooking{}, err
	}

	var row = &entity.CarBooking{}
	err = attributevalue.UnmarshalMap(res.Item, row)
	if err != nil {
		return dto.CarBooking{}, err
	}

	return row.DTO(), nil
}
