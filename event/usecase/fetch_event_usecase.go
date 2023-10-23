package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fetchEventUsecase struct {
	fetchEventRepo domain.EventFetchRepository
	contextTimeout time.Duration
}

func (eventUsecase *fetchEventUsecase) InsertOne(c context.Context, e *domain.EventFetch) (*domain.EventFetch, error) {
	ctx, cancel := context.WithTimeout(c, eventUsecase.contextTimeout)
	defer cancel()

	e.ID = primitive.NewObjectID()
	e.LastFetch = time.Now()

	res, err := eventUsecase.fetchEventRepo.InsertOne(ctx, e)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e fetchEventUsecase) FindOne(ctx context.Context, id string) (*domain.EventFetch, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.fetchEventRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e fetchEventUsecase) FindOneByType(ctx context.Context, typeOfEvent string) (*domain.EventFetch, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.fetchEventRepo.FindOneByType(ctx, typeOfEvent)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e fetchEventUsecase) UpdateOne(ctx context.Context, user *domain.EventFetch, id string) (*domain.EventFetch, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.fetchEventRepo.UpdateOne(ctx, user, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewFetchEentUsecase(u domain.EventFetchRepository, to time.Duration) domain.EventFetchUsecase {
	return &fetchEventUsecase{
		fetchEventRepo: u,
		contextTimeout: to,
	}

}
