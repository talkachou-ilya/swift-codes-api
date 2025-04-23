package test_cases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"swift-codes-api/models"
)

type IntegrationTestCase struct {
	Name               string
	SwiftCode          string
	SetupData          func(context.Context, *mongo.Database)
	ExpectedStatusCode int
	ExpectedResponse   string
}

func GetIntegrationTestCases() []IntegrationTestCase {
	return []IntegrationTestCase{
		{
			Name:      "Headquarter SWIFT code with branches",
			SwiftCode: "CITIUS33XXX",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				hq := models.SwiftCode{
					SwiftCode:     "CITIUS33XXX",
					SwiftPrefix:   "CITIUS",
					IsHeadquarter: true,
					BankName:      "Citibank",
					Address:       "388 Greenwich Street, New York",
					CountryISO2:   "US",
					CountryName:   "United States",
				}
				branch := models.SwiftCode{
					SwiftCode:     "CITIUS22XXX",
					SwiftPrefix:   "CITIUS",
					IsHeadquarter: false,
					BankName:      "Citibank",
					Address:       "1 Court Square, Long Island City",
					CountryISO2:   "US",
					CountryName:   "United States",
				}

				insertResult, insertErr := db.Collection("swift-codes").InsertMany(ctx, []interface{}{hq, branch})
				log.Printf("Insert result: %+v, Error: %v", insertResult, insertErr)
				documents, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Println("Documents in collection:", documents)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"address": "388 Greenwich Street, New York",
				"bankName": "Citibank",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": true,
				"swiftCode": "CITIUS33XXX",
				"branches": [
					{
						"address": "1 Court Square, Long Island City",
						"bankName": "Citibank",
						"countryISO2": "US",
						"countryName": "United States",
						"isHeadquarter": false,
						"swiftCode": "CITIUS22XXX"
					}
				]
			}`,
		},
	}
}
