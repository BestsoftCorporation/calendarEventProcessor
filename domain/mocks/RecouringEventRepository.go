package mocks

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/stretchr/testify/mock"
)

type RecouringEvenetRepository struct {
	mock.Mock
}

// InsertOne provides a mock function with given fields: ctx, u
func (_m *RecouringEvenetRepository) InsertOne(ctx context.Context, u *domain.StandardTrip) (*domain.StandardTrip, error) {
	ret := _m.Called(ctx, u)

	var r0 *domain.StandardTrip
	if rf, ok := ret.Get(0).(func(context.Context, *domain.StandardTrip) *domain.StandardTrip); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.StandardTrip)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.StandardTrip) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
