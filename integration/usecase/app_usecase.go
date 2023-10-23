package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type appUsecase struct {
	appRepo        domain.AppRepository
	contextTimeout time.Duration
}

func NewAppUsecase(u domain.AppRepository, to time.Duration) domain.AppUsecase {
	return &appUsecase{
		appRepo:        u,
		contextTimeout: to,
	}
}

func (a appUsecase) InsertOne(ctx context.Context, app *domain.App) (*domain.App, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	app.ID = primitive.NewObjectID()

	res, err := a.appRepo.InsertOne(ctx, app)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (a appUsecase) FindOne(ctx context.Context, token string) (*domain.App, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res, err := a.appRepo.FindOne(ctx, token)
	if err != nil {
		return res, err
	}

	return res, nil
}
