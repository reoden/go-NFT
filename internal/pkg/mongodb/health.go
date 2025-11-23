package mongodb

import (
	"context"

	"github.com/reoden/go-echo-template/pkg/health/contracts"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoHealthChecker struct {
	client *mongo.Client
}

func NewMongoHealthChecker(client *mongo.Client) contracts.Health {
	return &mongoHealthChecker{client}
}

func (healthChecker *mongoHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx, nil)
}

func (healthChecker *mongoHealthChecker) GetHealthName() string {
	return "mongodb"
}
