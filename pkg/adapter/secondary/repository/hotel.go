package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant/status"
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/entity"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type HotelRepository struct {
	client *dynamodb.Client
}

func NewHotelRepository(client *dynamodb.Client) *HotelRepository {
	return &HotelRepository{
		client: client,
	}
}

func (r *HotelRepository) Book(ctx context.Context, tripID string) (dto.HotelBooking, error) {
	if booking, err := r.GetByTripID(ctx, tripID); err != nil {
		return dto.HotelBooking{}, err
	} else if booking.Status == status.Booked {
		return booking, nil
	}

	condition := expression.AttributeNotExists(expression.Name("PK"))
	condition.And(expression.AttributeNotExists(expression.Name("SK")))
	condExpr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return dto.HotelBooking{}, err
	}

	req := &entity.HotelBooking{
		PK:        fmt.Sprintf("BOOKING#HOTEL#%s", tripID),
		SK:        fmt.Sprintf("#BOOKING#HOTEL#%s", tripID),
		ID:        uuid.NewString(),
		TripID:    tripID,
		HotelID:   uuid.NewString(),
		Status:    status.Booked,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if err := util.ValidateStruct(req); err != nil {
		return dto.HotelBooking{}, err
	}

	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		return dto.HotelBooking{}, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("trip"),
		Item:                item,
		ConditionExpression: condExpr.Condition(),
	})
	if err != nil {
		return dto.HotelBooking{}, err
	}

	return req.DTO(), nil
}

func (r *HotelRepository) GetByTripID(ctx context.Context, tripID string) (dto.HotelBooking, error) {
	pk, err := attributevalue.Marshal(fmt.Sprintf("BOOKING#HOTEL#%s", tripID))
	if err != nil {
		return dto.HotelBooking{}, err
	}
	sk, err := attributevalue.Marshal(fmt.Sprintf("#BOOKING#HOTEL#%s", tripID))
	if err != nil {
		return dto.HotelBooking{}, err
	}

	res, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("trip"),
		Key:       map[string]types.AttributeValue{"PK": pk, "SK": sk},
	})
	if err != nil {
		return dto.HotelBooking{}, err
	}

	var row = &entity.HotelBooking{}
	err = attributevalue.UnmarshalMap(res.Item, row)
	if err != nil {
		return dto.HotelBooking{}, err
	}

	return row.DTO(), nil
}
