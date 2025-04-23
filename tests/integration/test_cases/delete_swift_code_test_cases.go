package test_cases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"swift-codes-api/models"
)

type DeleteSwiftCodeIntegrationTestCase struct {
	Name               string
	SwiftCode          string
	SetupData          func(context.Context, *mongo.Database)
	CheckData          func(context.Context, *mongo.Database) bool
	ExpectedStatusCode int
	ExpectedResponse   string
}

func GetDeleteSwiftCodeIntegrationTestCases() []DeleteSwiftCodeIntegrationTestCase {
	return []DeleteSwiftCodeIntegrationTestCase{
		{
			Name:      "Successful deletion",
			SwiftCode: "DELEUS33",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				swiftCode := models.SwiftCode{
					SwiftCode:     "DELEUS33",
					SwiftPrefix:   "DELEUS",
					BankName:      "DELETE TEST BANK",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 Delete St, New York",
					IsHeadquarter: true,
				}
				result, err := db.Collection("swift-codes").InsertOne(ctx, swiftCode)
				log.Printf("Insert result: %+v, Error: %v", result, err)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "DELEUS33"})
				return err == nil && count == 0
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"message":"SWIFT code deleted successfully"}`,
		},
		{
			Name:      "Swift code not found",
			SwiftCode: "NOTFOUND",
			SetupData: func(ctx context.Context, db *mongo.Database) {
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				return true
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"SWIFT code not found"}`,
		},
		{
			Name:      "Delete one of multiple codes",
			SwiftCode: "DELETE01",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				codes := []interface{}{
					models.SwiftCode{
						SwiftCode:     "DELETE01",
						SwiftPrefix:   "DELETE",
						BankName:      "DELETE BANK ONE",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "1 Delete Blvd, New York",
						IsHeadquarter: true,
					},
					models.SwiftCode{
						SwiftCode:     "DELETE02",
						SwiftPrefix:   "DELETE",
						BankName:      "DELETE BANK TWO",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "2 Delete Blvd, New York",
						IsHeadquarter: false,
					},
				}
				result, err := db.Collection("swift-codes").InsertMany(ctx, codes)
				log.Printf("Insert result: %+v, Error: %v", result, err)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count1, err1 := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "DELETE01"})
				count2, err2 := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "DELETE02"})
				return err1 == nil && count1 == 0 && err2 == nil && count2 == 1
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"message":"SWIFT code deleted successfully"}`,
		},
	}
}
