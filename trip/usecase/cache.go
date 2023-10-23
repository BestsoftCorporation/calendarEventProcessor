package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type cacheUsecase struct {
	cacheRepo      domain.CacheRepository
	contextTimeout time.Duration
}

func NewCacheUsecase(u domain.CacheRepository, to time.Duration) domain.CacheUsecase {
	return &cacheUsecase{
		cacheRepo:      u,
		contextTimeout: to,
	}

}

func (c cacheUsecase) InsertOne(ctx context.Context, cache *domain.Cache) (*domain.Cache, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	cache.ID = primitive.NewObjectID()
	res, err := c.cacheRepo.InsertOne(ctx, cache)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c cacheUsecase) FindOne(ctx context.Context, startLocation string, endLocation string, mode string, day string) (*domain.Cache, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	res, err := c.cacheRepo.FindOne(ctx, startLocation, endLocation, mode, day)
	if err != nil {
		return res, err
	}

	return res, nil
}
