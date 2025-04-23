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
			Name:      "Successful deletion of headquarter SWIFT code",
			SwiftCode: "BOTKUS33XXX",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				swiftCode := models.SwiftCode{
					SwiftCode:     "BOTKUS33XXX",
					SwiftPrefix:   "BOTKUS",
					BankName:      "Bank of Test",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 Test Avenue, New York",
					IsHeadquarter: true,
				}
				result, err := db.Collection("swift-codes").InsertOne(ctx, swiftCode)
				log.Printf("Insert result: %+v, Error: %v", result, err)

				count, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Printf("Documents before test: %d", count)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "BOTKUS33XXX"})
				if err != nil {
					log.Printf("Error checking data: %v", err)
					return false
				}
				log.Printf("Documents after deletion: %d", count)
				return count == 0
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"message":"SWIFT code deleted successfully"}`,
		},
		{
			Name:      "SWIFT code not found",
			SwiftCode: "NOTFND33XXX",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				swiftCode := models.SwiftCode{
					SwiftCode:     "TSTUKB1XXX",
					SwiftPrefix:   "TSTUKB",
					BankName:      "Test Bank UK",
					CountryISO2:   "GB",
					CountryName:   "United Kingdom",
					Address:       "1 Test Street, London",
					IsHeadquarter: true,
				}
				result, err := db.Collection("swift-codes").InsertOne(ctx, swiftCode)
				log.Printf("Insert result: %+v, Error: %v", result, err)

				count, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Printf("Documents before test: %d", count)
			},
			CheckData: func(ctx context.Context, db *mongo.Database) bool {
				count, err := db.Collection("swift-codes").CountDocuments(ctx, bson.M{"swiftCode": "TSTUKB1XXX"})
				if err != nil {
					log.Printf("Error checking data: %v", err)
					return false
				}
				log.Printf("Documents after test: %d", count)
				return count == 1
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"SWIFT code not found"}`,
		},
	}
}
