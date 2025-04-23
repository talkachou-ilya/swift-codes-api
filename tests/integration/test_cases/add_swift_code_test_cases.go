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
			Name: "Valid headquarter SWIFT code",
			RequestBody: `{
				"swiftCode": "TESTUS33XXX",
				"bankName": "Test Bank America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Test Avenue, New York",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				count, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Printf("Documents before test: %d", count)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				var result models.SwiftCode
				err := db.Collection("swift-codes").FindOne(ctx, bson.M{"swiftCode": "TESTUS33XXX"}).Decode(&result)
				if err != nil {
					log.Printf("Error checking data: %v", err)
					return false
				}

				log.Printf("Found SWIFT code: %s, BankName: %s, SwiftPrefix: %s",
					result.SwiftCode, result.BankName, result.SwiftPrefix)

				return result.BankName == "Test Bank America" &&
					result.SwiftPrefix == "TESTUS33" &&
					result.IsHeadquarter == true
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse:   `{"message":"SWIFT code added successfully"}`,
		},
		{
			Name: "Duplicate SWIFT code",
			RequestBody: `{
				"swiftCode": "DUPEFR33XXX",
				"bankName": "Duplicate Bank",
				"countryISO2": "FR",
				"countryName": "France",
				"address": "123 Duplicate Street, Paris",
				"isHeadquarter": true
			}`,
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				existingCode := models.SwiftCode{
					SwiftCode:     "DUPEFR33XXX",
					SwiftPrefix:   "DUPEFR33",
					BankName:      "Existing Bank France",
					CountryISO2:   "FR",
					CountryName:   "France",
					Address:       "456 Existing Blvd, Paris",
					IsHeadquarter: true,
				}
				insertResult, err := db.Collection("swift-codes").InsertOne(ctx, existingCode)
				log.Printf("Insert result: %+v, Error: %v", insertResult, err)

				count, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Printf("Documents before test: %d", count)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				var result models.SwiftCode
				err := db.Collection("swift-codes").FindOne(ctx, bson.M{"swiftCode": "DUPEFR33XXX"}).Decode(&result)
				if err != nil {
					log.Printf("Error checking data: %v", err)
					return false
				}

				log.Printf("Found SWIFT code: %s, BankName: %s", result.SwiftCode, result.BankName)
				return result.BankName == "Existing Bank France"
			},
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   `{"message":"SWIFT code DUPEFR33XXX already exists"}`,
		},
	}
}
