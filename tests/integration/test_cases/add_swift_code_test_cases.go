package test_cases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"swift-codes-api/models"
)

type AddSwiftCodeIntegrationTestCase struct {
	Name               string
	RequestBody        string
	SetupData          func(context.Context, *mongo.Database)
	CheckData          func(context.Context, *mongo.Database) bool
	ExpectedStatusCode int
	ExpectedResponse   string
}

func GetAddSwiftCodeIntegrationTestCases() []AddSwiftCodeIntegrationTestCase {
	return []AddSwiftCodeIntegrationTestCase{
		{
			Name: "Valid SWIFT code",
			RequestBody: `{
				"swiftCode": "TESTUS33",
				"bankName": "Test Bank",
				"countryISO2": "us",
				"countryName": "United States",
				"address": "123 Test St, New York",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				var result models.SwiftCode
				err := db.Collection("swift-codes").FindOne(ctx, bson.M{"swiftCode": "TESTUS33"}).Decode(&result)
				return err == nil && result.BankName == "Test Bank" && result.SwiftPrefix == "TESTUS"
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse:   `{"message":"SWIFT code added successfully"}`,
		},
		{
			Name: "Duplicate SWIFT code",
			RequestBody: `{
				"swiftCode": "DUPEUS33",
				"bankName": "Duplicate Bank",
				"countryISO2": "us",
				"countryName": "United States",
				"address": "123 Dupe St, New York",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
				existingCode := models.SwiftCode{
					SwiftCode:     "DUPEUS33",
					SwiftPrefix:   "DUPEUS",
					BankName:      "Existing Bank",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "456 Existing St, Chicago",
					IsHeadquarter: true,
				}
				insertResult, err := db.Collection("swift-codes").InsertOne(ctx, existingCode)
				log.Printf("Insert result: %+v, Error: %v", insertResult, err)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				var result models.SwiftCode
				err := db.Collection("swift-codes").FindOne(ctx, bson.M{"swiftCode": "DUPEUS33"}).Decode(&result)
				return err == nil && result.BankName == "Existing Bank"
			},
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   `{"message":"SWIFT code DUPEUS33 already exists"}`,
		},
		{
			Name: "Invalid SWIFT code format",
			RequestBody: `{
				"swiftCode": "SHORT",
				"bankName": "Short Code Bank",
				"countryISO2": "us",
				"countryName": "United States",
				"address": "123 Short St, New York",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "SHORT"})
				return err == nil && count == 0
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"message":"Invalid SWIFT code format. Must be 8 or 11 characters"}`,
		},
		{
			Name: "Invalid country code format",
			RequestBody: `{
				"swiftCode": "TESTUS33",
				"bankName": "Test Bank",
				"countryISO2": "USA",
				"countryName": "United States",
				"address": "123 Test St, New York",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "TESTUS33"})
				return err == nil && count == 0
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"message":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name: "Missing required fields",
			RequestBody: `{
				"swiftCode": "TESTUS33",
				"bankName": "",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {

			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {

				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "TESTUS33"})
				return err == nil && count == 0
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"message":"Missing required fields"}`,
		},
	}
}
