package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventUsecase struct {
	eventRepo      domain.EventRepository
	contextTimeout time.Duration
}

func (eventUsecase *eventUsecase) InsertOne(c context.Context, e *domain.Event) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(c, eventUsecase.contextTimeout)
	defer cancel()

	e.ID = primitive.NewObjectID()
	e.CreatedDate = time.Now()

	res, err := eventUsecase.eventRepo.InsertOne(ctx, e)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) FindOne(ctx context.Context, id string) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) Find(ctx context.Context, email string, date string) (*[]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.Find(ctx, email, date)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) FindCommute(ctx context.Context, date string) (*[]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.FindCommute(ctx, date)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) UpdateOne(ctx context.Context, user *domain.Event, id string) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.UpdateOne(ctx, user, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) FindOneById(ctx context.Context, id string) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.FindOneByID(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) DeleteOneByID(ctx context.Context, id string) (*domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.DeleteOneByID(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (e eventUsecase) DeleteAll(ctx context.Context, user_email string) (*[]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()

	res, err := e.eventRepo.DeleteAll(ctx, user_email)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewEventUsecase(u domain.EventRepository, to time.Duration) domain.EventUsecase {
	return &eventUsecase{
		eventRepo:      u,
		contextTimeout: to,
	}

}
