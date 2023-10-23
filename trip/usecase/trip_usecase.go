package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripUsecase struct {
	TripRepo       domain.TripRepository
	contextTimeout time.Duration
}

func NewTripUsecase(u domain.TripRepository, to time.Duration) domain.TripUsecase {
	return &TripUsecase{
		TripRepo:       u,
		contextTimeout: to,
	}

}

func (tripUsecase *TripUsecase) InsertOne(c context.Context, e *domain.Trip) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(c, tripUsecase.contextTimeout)
	defer cancel()

	e.ID = primitive.NewObjectID()
	e.CreatedDate = time.Now()

	res, err := tripUsecase.TripRepo.InsertOne(ctx, e)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e TripUsecase) FindOne(ctx context.Context, event_id string, returning bool) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.TripRepo.FindOne(ctx, event_id, returning)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e TripUsecase) DeleteAll(ctx context.Context, userEmail string) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.TripRepo.DeleteAll(ctx, userEmail)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e TripUsecase) UpdateOne(ctx context.Context, trip *domain.Trip, id primitive.ObjectID) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.TripRepo.UpdateOne(ctx, trip, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (tripUsecase *TripUsecase) FindAll(ctx context.Context, userEmail string) (*[]domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, tripUsecase.contextTimeout)
	defer cancel()

	res, err := tripUsecase.TripRepo.FindAll(ctx, userEmail)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e TripUsecase) DeleteByReturningAndID(ctx context.Context, ret bool, id string, field string, disableID string,
	delete bool) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.TripRepo.DeleteByReturningAndID(ctx, ret, id, field, disableID, delete)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e TripUsecase) DeleteAllWhereLinkedId(ctx context.Context, linkedId string) (*domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.TripRepo.DeleteAllWhereLinkedId(ctx, linkedId)
	if err != nil {
		return res, err
	}

	return res, nil
}
