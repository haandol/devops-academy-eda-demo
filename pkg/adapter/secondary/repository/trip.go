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
	"github.com/haandol/devops-academy-eda-demo/pkg/constant/status"
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/entity"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type TripRepository struct {
	client *dynamodb.Client
}

func NewTripRepository(client *dynamodb.Client) *TripRepository {
	return &TripRepository{
		client: client,
	}
}

func (r *TripRepository) Create(ctx context.Context, tripID string) (dto.Trip, error) {
	condition := expression.AttributeNotExists(expression.Name("PK"))
	condition.And(expression.AttributeNotExists(expression.Name("SK")))
	condExpr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return dto.Trip{}, err
	}

	now := time.Now().Format(time.RFC3339)
	req := &entity.Trip{
		PK:        fmt.Sprintf("TRIP#%s", tripID),
		SK:        fmt.Sprintf("#TRIP#%s", tripID),
		GS1PK:     fmt.Sprintf("STATUS#%s", status.TripInitialized),
		GS1SK:     fmt.Sprintf("#CREATED#%s", now),
		ID:        tripID,
		Status:    status.TripInitialized,
		CreatedAt: now,
	}
	if err := util.ValidateStruct(req); err != nil {
		return dto.Trip{}, err
	}

	item, err := attributevalue.MarshalMap(req)
	if err != nil {
		return dto.Trip{}, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String("trip"),
		Item:                item,
		ConditionExpression: condExpr.Condition(),
	})
	if err != nil {
		return dto.Trip{}, err
	}

	return req.DTO(), nil
}

func (r *TripRepository) Complete(ctx context.Context, evt *event.FlightBooked) error {
	now := time.Now().Format(time.RFC3339)
	update := expression.Set(expression.Name("GS1PK"), expression.Value(fmt.Sprintf("STATUS#%s", status.TripBooked)))
	update.Set(expression.Name("flightID"), expression.Value(evt.Body.FlightID))
	update.Set(expression.Name("flightBookingID"), expression.Value(evt.Body.BookingID))
	update.Set(expression.Name("status"), expression.Value(status.TripBooked))
	update.Set(expression.Name("updatedAt"), expression.Value(now))
	updateExpr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	condition := expression.AttributeExists(expression.Name("PK"))
	condition.And(expression.AttributeExists(expression.Name("SK")))
	condExpr, err := expression.NewBuilder().WithCondition(condition).Build()
	if err != nil {
		return err
	}

	pk, err := attributevalue.Marshal(fmt.Sprintf("TRIP#%s", evt.Body.TripID))
	if err != nil {
		return err
	}
	sk, err := attributevalue.Marshal(fmt.Sprintf("#TRIP#%s", evt.Body.TripID))
	if err != nil {
		return err
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String("trip"),
		Key:                       map[string]types.AttributeValue{"PK": pk, "SK": sk},
		ExpressionAttributeNames:  updateExpr.Names(),
		ExpressionAttributeValues: updateExpr.Values(),
		UpdateExpression:          updateExpr.Update(),
		ConditionExpression:       condExpr.Condition(),
	})
	return err
}

func (r *TripRepository) List(ctx context.Context) ([]dto.Trip, error) {
	var rows []dto.Trip

	return rows, nil
}
