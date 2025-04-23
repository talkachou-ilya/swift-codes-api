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

func (m *SwiftRepository) FindByCountryISO2(ctx context.Context, countryISO2 string) ([]models.SwiftCode, string, error) {
	args := m.Called(ctx, countryISO2)
	if args.Get(0) != nil && args.Get(1) != nil {
		return args.Get(0).([]models.SwiftCode), args.String(1), args.Error(2)
	}
	return nil, "", args.Error(2)
}

func (m *SwiftRepository) AddSwiftCode(ctx context.Context, swiftCode models.SwiftCode) error {
	args := m.Called(ctx, swiftCode)
	return args.Error(0)
}

func (m *SwiftRepository) DeleteSwiftCode(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}
