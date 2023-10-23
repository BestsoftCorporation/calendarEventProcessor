package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	_eventCalendar "github.com/bxcodec/go-clean-arch/event/delivery/calendar"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type recurringEventUsecase struct {
	recurringEventRepo domain.StandardTripRepository
	contextTimeout     time.Duration
}

func NewRecuringEventUsecase(u domain.StandardTripRepository, to time.Duration) domain.StandardTripUsecase {
	return &recurringEventUsecase{
		recurringEventRepo: u,
		contextTimeout:     to,
	}

}

func (recuringEventUsecase *recurringEventUsecase) InsertOne(c context.Context, m *domain.StandardTrip) (*domain.StandardTrip, error) {

	ctx, cancel := context.WithTimeout(c, recuringEventUsecase.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := recuringEventUsecase.recurringEventRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

	//fillDistances(*m)
	for _, day := range m.Days {
		_eventCalendar.PushEvents(m.UserEmail, day.Mode, day.Day, m.UserWorkspaceAddress)
	}

	return res, nil
}

func (recuringEventUsecase recurringEventUsecase) FindOne(ctx context.Context, id string) (*domain.StandardTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, recuringEventUsecase.contextTimeout)
	defer cancel()

	res, err := recuringEventUsecase.recurringEventRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (recuringEventUsecase *recurringEventUsecase) FindAll(ctx context.Context) (*[]domain.StandardTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, recuringEventUsecase.contextTimeout)
	defer cancel()

	res, err := recuringEventUsecase.recurringEventRepo.FindAll(ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (recuringEventUsecase *recurringEventUsecase) DeleteOne(ctx context.Context, email string) (*domain.StandardTrip, error) {
	ctx, cancel := context.WithTimeout(ctx, recuringEventUsecase.contextTimeout)
	defer cancel()

	res, err := recuringEventUsecase.recurringEventRepo.DeleteOne(ctx, email)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (recuringEventUsecase recurringEventUsecase) UpdateOne(ctx context.Context, user *domain.StandardTrip, id string) (*domain.StandardTrip, error) {
	//TODO implement me
	panic("implement me")
}
