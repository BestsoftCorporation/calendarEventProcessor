package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventFetch struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	LastFetch time.Time          `bson:"last_fetch" json:"last_fetch"`
	Type      string             `bson:"type" json:"type"`
	Processed bool               `bson:"processed" json:"processed"`
}

type EventFetchRepository interface {
	InsertOne(ctx context.Context, fetchEvent *EventFetch) (*EventFetch, error)
	FindOne(ctx context.Context, id string) (*EventFetch, error)
	FindOneByType(ctx context.Context, typeOfEvent string) (*EventFetch, error)
	UpdateOne(ctx context.Context, fetchEvent *EventFetch, id string) (*EventFetch, error)
}

type EventFetchUsecase interface {
	InsertOne(ctx context.Context, fetchEvent *EventFetch) (*EventFetch, error)
	FindOne(ctx context.Context, id string) (*EventFetch, error)
	FindOneByType(ctx context.Context, typeOfEvent string) (*EventFetch, error)
	UpdateOne(ctx context.Context, user *EventFetch, id string) (*EventFetch, error)
}
