package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"swift-codes-api/models"
)

type SwiftRepository struct {
	mock.Mock
}

func (m *SwiftRepository) FindByCode(ctx context.Context, code string) (*models.SwiftCode, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*models.SwiftCode), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SwiftRepository) FindBranchesByPrefix(ctx context.Context, prefix string) ([]models.SwiftCode, error) {
	args := m.Called(ctx, prefix)
	if args.Get(0) != nil {
		return args.Get(0).([]models.SwiftCode), args.Error(1)
	}
	return nil, args.Error(1)
}
