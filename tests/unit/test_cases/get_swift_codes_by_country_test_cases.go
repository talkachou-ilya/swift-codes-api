package test_cases

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"swift-codes-api/models"
	mockRepo "swift-codes-api/repositories/mock"
)

type CountryTestCase struct {
	Name             string
	CountryISO2      string
	SetupMocks       func(repository *mockRepo.SwiftRepository)
	ExpectedStatus   int
	ExpectedResponse string
}

func GetCountryTestCases() []CountryTestCase {
	return []CountryTestCase{
		{
			Name:        "valid country code",
			CountryISO2: "US",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCodes := []models.SwiftCode{
					{
						SwiftCode:     "AAAAUS33XXX",
						BankName:      "BANK OF AMERICA",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "123 MAIN ST, NEW YORK",
						IsHeadquarter: true,
					},
					{
						SwiftCode:     "BBBBUS44XXX",
						BankName:      "CHASE BANK",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "456 OAK ST, CHICAGO",
						IsHeadquarter: true,
					},
				}
				repo.On("FindByCountryISO2", mock.Anything, "US").Return(
					swiftCodes,
					"United States",
					nil,
				)
			},
			ExpectedStatus: 200,
			ExpectedResponse: `{
				"countryISO2": "US",
				"countryName": "United States",
				"swiftCodes": [
					{
						"swiftCode": "AAAAUS33XXX",
						"bankName": "BANK OF AMERICA",
						"countryISO2": "US",
						"countryName": "United States",
						"address": "123 MAIN ST, NEW YORK",
						"isHeadquarter": true
					},
					{
						"swiftCode": "BBBBUS44XXX",
						"bankName": "CHASE BANK",
						"countryISO2": "US",
						"countryName": "United States",
						"address": "456 OAK ST, CHICAGO",
						"isHeadquarter": true
					}
				]
			}`,
		},
		{
			Name:             "invalid country code length",
			CountryISO2:      "USA",
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   400,
			ExpectedResponse: `{"error":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name:        "country not found",
			CountryISO2: "FR",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("FindByCountryISO2", mock.Anything, "FR").Return(
					[]models.SwiftCode{},
					"",
					mongo.ErrNoDocuments,
				)
			},
			ExpectedStatus:   404,
			ExpectedResponse: `{"message":"No SWIFT codes found for this country"}`,
		},
	}
}
