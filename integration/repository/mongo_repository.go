package repository

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAppRepository struct {
	DB         mongo.Database
	Collection *mongo.Collection
}

const (
	cacheCollectionName = "apps"
)

func NewMongoAppRepository(DB mongo.Database) domain.AppRepository {
	return &mongoAppRepository{DB, DB.Collection(cacheCollectionName)}
}

func (m mongoAppRepository) InsertOne(ctx context.Context, app *domain.App) (*domain.App, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, app)

	if err != nil {
		return app, err
	}

	return app, nil
}

func (m mongoAppRepository) FindOne(ctx context.Context, Token string) (*domain.App, error) {
	var (
		app domain.App
		err error
	)

	err = m.Collection.FindOne(ctx, bson.M{"token": Token}).Decode(&app)
	if err != nil {
		return &app, err
	}

	return &app, nil
}
