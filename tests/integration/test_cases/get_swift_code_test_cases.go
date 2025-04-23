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
			Name:      "Valid headquarter SWIFT code with branch",
			SwiftCode: "TESTUSXXHQ",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				hq := models.SwiftCode{
					SwiftCode:     "TESTUSXXHQ",
					SwiftPrefix:   "TESTUSXX",
					IsHeadquarter: true,
					BankName:      "Test HQ Bank",
					Address:       "HQ Street",
					CountryISO2:   "US",
					CountryName:   "United States",
				}
				branch := models.SwiftCode{
					SwiftCode:     "TESTUSXX01",
					SwiftPrefix:   "TESTUSXX",
					IsHeadquarter: false,
					BankName:      "Test HQ Bank",
					Address:       "Branch Street",
					CountryISO2:   "US",
					CountryName:   "United States",
				}

				insertResult, insertErr := db.Collection("swift-codes").InsertMany(ctx, []interface{}{hq, branch})
				log.Printf("Insert result: %+v, Error: %v", insertResult, insertErr)
				documents, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Println(documents)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"address": "HQ Street",
				"bankName": "Test HQ Bank",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": true,
				"swiftCode": "TESTUSXXHQ",
				"branches": [
					{
						"address": "Branch Street",
						"bankName": "Test HQ Bank",
						"countryISO2": "US",
						"countryName": "United States",
						"isHeadquarter": false,
						"swiftCode": "TESTUSXX01"
					}
				]
			}`,
		},
	}
}
