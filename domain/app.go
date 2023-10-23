package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	AppName string             `bson:"app_name" json:"app_name"`
	Token   string             `bson:"token" json:"token"`
}

type AppRepository interface {
	InsertOne(ctx context.Context, app *App) (*App, error)
	FindOne(ctx context.Context, token string) (*App, error)
}

type AppUsecase interface {
	InsertOne(ctx context.Context, cache *App) (*App, error)
	FindOne(ctx context.Context, token string) (*App, error)
}
