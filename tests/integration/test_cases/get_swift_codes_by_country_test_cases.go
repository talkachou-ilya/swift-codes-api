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
			CountryISO2: "DE",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				banks := []models.SwiftCode{
					{
						SwiftCode:     "DEUTDEFF",
						SwiftPrefix:   "DEUTDE",
						IsHeadquarter: true,
						BankName:      "Deutsche Bank",
						Address:       "Taunusanlage 12, Frankfurt am Main",
						CountryISO2:   "DE",
						CountryName:   "Germany",
					},
					{
						SwiftCode:     "COBADEFF",
						SwiftPrefix:   "COBADE",
						IsHeadquarter: true,
						BankName:      "Commerzbank",
						Address:       "Kaiserplatz, Frankfurt am Main",
						CountryISO2:   "DE",
						CountryName:   "Germany",
					},
				}

				insertResult, insertErr := db.Collection("swift-codes").InsertMany(ctx, []interface{}{banks[0], banks[1]})
				log.Printf("Insert result: %+v, Error: %v", insertResult, insertErr)
				documents, _ := db.Collection("swift-codes").CountDocuments(ctx, bson.M{})
				log.Println(documents)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: `{
				"countryISO2": "DE",
				"countryName": "Germany",
				"swiftCodes": [
					{
						"swiftCode": "DEUTDEFF",
						"bankName": "Deutsche Bank",
						"address": "Taunusanlage 12, Frankfurt am Main",
						"countryISO2": "DE",
						"countryName": "Germany",
						"isHeadquarter": true
					},
					{
						"swiftCode": "COBADEFF",
						"bankName": "Commerzbank",
						"address": "Kaiserplatz, Frankfurt am Main",
						"countryISO2": "DE",
						"countryName": "Germany",
						"isHeadquarter": true
					}
				]
			}`,
		},
		{
			Name:        "Country with no banks",
			CountryISO2: "ZZ",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				// No data to insert for this test case
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"No SWIFT codes found for this country"}`,
		},
		{
			Name:        "Invalid country code format",
			CountryISO2: "USA",
			SetupData: func(ctx context.Context, db *mongo.Database) {
				// No data to insert for this test case
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"error":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
	}
}
