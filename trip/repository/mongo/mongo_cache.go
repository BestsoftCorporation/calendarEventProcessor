package mongo

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCacheRepository struct {
	DB         mongo.Database
	Collection *mongo.Collection
}

const (
	cacheCollectionName = "cache"
)

func NewMongoCacheRepository(DB mongo.Database) domain.CacheRepository {
	return &mongoCacheRepository{DB, DB.Collection(cacheCollectionName)}
}

func (m mongoCacheRepository) InsertOne(ctx context.Context, cache *domain.Cache) (*domain.Cache, error) {
	var (
		err error
	)

	_, err = m.Collection.InsertOne(ctx, cache)

	if err != nil {
		return cache, err
	}

	return cache, nil
}

func (m mongoCacheRepository) FindOne(ctx context.Context, startLocation string, endLocation string, mode string, day string) (*domain.Cache, error) {
	var (
		cache domain.Cache
		err   error
	)

	err = m.Collection.FindOne(ctx, bson.M{"start_location": startLocation, "end_location": endLocation, "mode": mode, "day": day}).Decode(&cache)
	if err != nil {
		return &cache, err
	}

	return &cache, nil
}
