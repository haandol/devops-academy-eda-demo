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

type FlightRepository struct {
	client *dynamodb.Client
}

func NewFlightRepository(client *dynamodb.Client) *FlightRepository {
	return &FlightRepository{
		client: client,
	}
}

func (r *FlightRepository) Book(ctx context.Context, tripID string) (dto.FlightBooking, error) {
	if booking, err := r.GetByTripID(ctx, tripID); err != nil {
		return dto.FlightBooking{}, err
	} else if booking.Status == status.Booked {
		return booking, nil
	}

	condition := expression.AttributeNotExists(expression.Name("PK"))
	condition.And(expression.AttributeNotExists(expression.Name("SK")))
	condExpr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return dto.FlightBooking{}, err
	}

	req := &entity.FlightBooking{
		PK:        fmt.Sprintf("BOOKING#FLIGHT#%s", tripID),
		SK:        fmt.Sprintf("#BOOKING#FLIGHT#%s", tripID),
		ID:        uuid.NewString(),
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

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("trip"),
		Item:                item,
		ConditionExpression: condExpr.Condition(),
	})
	if err != nil {
		return dto.FlightBooking{}, err
	}

	return req.DTO(), nil
}

func (r *FlightRepository) GetByTripID(ctx context.Context, tripID string) (dto.FlightBooking, error) {
	pk, err := attributevalue.Marshal(fmt.Sprintf("BOOKING#FLIGHT#%s", tripID))
	if err != nil {
		return dto.FlightBooking{}, err
	}
	sk, err := attributevalue.Marshal(fmt.Sprintf("#BOOKING#FLIGHT#%s", tripID))
	if err != nil {
		return dto.FlightBooking{}, err
	}

	res, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("trip"),
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
