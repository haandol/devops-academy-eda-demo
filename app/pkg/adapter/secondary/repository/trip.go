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
	"github.com/haandol/devops-academy-eda-demo/pkg/config"
	"github.com/haandol/devops-academy-eda-demo/pkg/constant/status"
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
	"github.com/haandol/devops-academy-eda-demo/pkg/entity"
	"github.com/haandol/devops-academy-eda-demo/pkg/message/event"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
)

type TripRepository struct {
	TableName string
	Client    *dynamodb.Client
}

func NewTripRepository(
	cfg *config.Config,
	client *dynamodb.Client,
) *TripRepository {
	return &TripRepository{
		TableName: cfg.Database.TableName,
		Client:    client,
	}
}

func (r *TripRepository) Create(ctx context.Context, tripID string) (dto.Trip, error) {
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

	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(r.TableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return dto.Trip{}, err
	}

	return req.DTO(), nil
}

func (r *TripRepository) Complete(ctx context.Context, evt *event.FlightBooked) error {
	now := time.Now().Format(time.RFC3339)
	update := expression.Set(expression.Name("GS1PK"), expression.Value(fmt.Sprintf("STATUS#%s", status.TripReserved))).
		Set(expression.Name("Status"), expression.Value(status.TripReserved)).
		Set(expression.Name("UpdatedAt"), expression.Value(now))
	updateExpr, err := expression.NewBuilder().WithUpdate(update).Build()
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

	_, err = r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(r.TableName),
		Key:                       map[string]types.AttributeValue{"PK": pk, "SK": sk},
		ExpressionAttributeNames:  updateExpr.Names(),
		ExpressionAttributeValues: updateExpr.Values(),
		UpdateExpression:          updateExpr.Update(),
		ConditionExpression:       aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	})
	return err
}

func (r *TripRepository) List(ctx context.Context) ([]dto.Trip, error) {
	filt := expression.Name("SK").BeginsWith("#TRIP#")
	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return []dto.Trip{}, err
	}

	res, err := r.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 aws.String(r.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		Limit:                     aws.Int32(100),
	})
	if err != nil {
		return []dto.Trip{}, err
	}

	var rows []dto.Trip
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &rows); err != nil {
		return []dto.Trip{}, err
	}

	return rows, nil
}
