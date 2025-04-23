package interfaces

import (
	"context"
	"swift-codes-api/models"
)

type SwiftRepository interface {
	FindByCode(ctx context.Context, code string) (*models.SwiftCode, error)
	FindBranchesByPrefix(ctx context.Context, prefix string) ([]models.SwiftCode, error)
	FindByCountryISO2(ctx context.Context, countryISO2 string) ([]models.SwiftCode, string, error)
}
