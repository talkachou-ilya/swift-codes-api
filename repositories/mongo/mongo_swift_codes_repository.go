package mongo

import (
	"context"
	"swift-codes-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SwiftRepository struct {
	col *mongo.Collection
}

func NewSwiftRepository(db *mongo.Database) *SwiftRepository {
	return &SwiftRepository{
		col: db.Collection("swift-codes"),
	}
}

func (r *SwiftRepository) FindByCode(ctx context.Context, code string) (*models.SwiftCode, error) {
	var result models.SwiftCode
	err := r.col.FindOne(ctx, bson.M{"swiftCode": code}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SwiftRepository) FindBranchesByPrefix(ctx context.Context, prefix string) ([]models.SwiftCode, error) {
	cursor, err := r.col.Find(ctx, bson.M{
		"swiftPrefix":   prefix,
		"isHeadquarter": false,
	})
	if err != nil {
		return nil, err
	}
	var branches []models.SwiftCode
	err = cursor.All(ctx, &branches)
	return branches, err
}

func (r *SwiftRepository) FindByCountryISO2(ctx context.Context, countryISO2 string) ([]models.SwiftCode, string, error) {
	cursor, err := r.col.Find(ctx, bson.M{"countryISO2": countryISO2})

	if err != nil {
		return nil, "", err
	}

	var swiftCodes []models.SwiftCode
	err = cursor.All(ctx, &swiftCodes)

	if len(swiftCodes) == 0 {
		return nil, "", mongo.ErrNoDocuments
	}

	countryName := swiftCodes[0].CountryName

	return swiftCodes, countryName, err
}
