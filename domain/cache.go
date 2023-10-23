package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cache struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	StartLocation string             `bson:"start_location" json:"start_location"`
	EndLocation   string             `bson:"end_location" json:"end_location"`
	Distance      int                `bson:"distance" json:"distance"`
	Mode          string             `bson:"mode" json:"mode"`
	Day           string             `bson:"day" json:"day"`
}

type CacheRepository interface {
	InsertOne(ctx context.Context, cache *Cache) (*Cache, error)
	FindOne(ctx context.Context, startLocation string, endLocation string, mode string, day string) (*Cache, error)
}

type CacheUsecase interface {
	InsertOne(ctx context.Context, cache *Cache) (*Cache, error)
	FindOne(ctx context.Context, startLocation string, endLocation string, mode string, day string) (*Cache, error)
}
