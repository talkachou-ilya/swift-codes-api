package test_cases

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"swift-codes-api/models"
)

type CountryIntegrationTestCase struct {
	Name               string
	CountryISO2        string
	SetupData          func(context.Context, *mongo.Database)
	ExpectedStatusCode int
	ExpectedResponse   string
}

func GetSwiftCodesByCountryIntegrationTestCases() []CountryIntegrationTestCase {
	return []CountryIntegrationTestCase{
		{
			Name:        "Valid country with multiple banks",
			CountryISO2: "US",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}

				banks := []models.SwiftCode{
					{
						SwiftCode:     "CITIUS33XXX",
						SwiftPrefix:   "CITIUS",
						IsHeadquarter: true,
						BankName:      "Citibank",
						Address:       "388 Greenwich Street, New York",
						CountryISO2:   "US",
						CountryName:   "United States",
					},
					{
						SwiftCode:     "CHASUS33XXX",
						SwiftPrefix:   "CHASUS",
						IsHeadquarter: true,
						BankName:      "JP Morgan Chase Bank",
						Address:       "270 Park Avenue, New York",
						CountryISO2:   "US",
						CountryName:   "United States",
					},
				}

				insertResult, insertErr := db.Collection("swift-codes").InsertMany(ctx, []interface{}{banks[0], banks[1]})
				log.Printf("Insert result: %+v, Error: %v", insertResult, insertErr)
				documents, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Println("Documents in collection:", documents)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"countryISO2": "US",
				"countryName": "United States",
				"swiftCodes": [
					{
						"swiftCode": "CITIUS33XXX",
						"bankName": "Citibank",
						"address": "388 Greenwich Street, New York",
						"countryISO2": "US",
						"countryName": "United States",
						"isHeadquarter": true
					},
					{
						"swiftCode": "CHASUS33XXX",
						"bankName": "JP Morgan Chase Bank",
						"address": "270 Park Avenue, New York",
						"countryISO2": "US",
						"countryName": "United States",
						"isHeadquarter": true
					}
				]
			}`,
		},
		{
			Name:        "Country with no banks",
			CountryISO2: "ZZ",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				_, err := db.Collection("swift-codes").DeleteMany(ctx, bson.M{})
				if err != nil {
					log.Printf("Error clearing collection: %v", err)
					return
				}
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"No SWIFT codes found for this country"}`,
		},
	}
}
