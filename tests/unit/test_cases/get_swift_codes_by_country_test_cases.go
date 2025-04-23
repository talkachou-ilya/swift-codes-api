package test_cases

import (
	"net/http"

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
			Name:        "Valid country code with multiple banks",
			CountryISO2: "US",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCodes := []models.SwiftCode{
					{
						SwiftCode:     "CITIUS33XXX",
						BankName:      "Citibank",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "388 Greenwich Street, New York",
						IsHeadquarter: true,
					},
					{
						SwiftCode:     "CHASUS33XXX",
						BankName:      "JP Morgan Chase Bank",
						CountryISO2:   "US",
						CountryName:   "United States",
						Address:       "270 Park Avenue, New York",
						IsHeadquarter: true,
					},
				}
				repo.On("FindByCountryISO2", mock.Anything, "US").Return(
					swiftCodes,
					"United States",
					nil,
				)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: `{
				"countryISO2": "US",
				"countryName": "United States",
				"swiftCodes": [
					{
						"swiftCode": "CITIUS33XXX",
						"bankName": "Citibank",
						"countryISO2": "US",
						"countryName": "United States",
						"address": "388 Greenwich Street, New York",
						"isHeadquarter": true
					},
					{
						"swiftCode": "CHASUS33XXX",
						"bankName": "JP Morgan Chase Bank",
						"countryISO2": "US",
						"countryName": "United States",
						"address": "270 Park Avenue, New York",
						"isHeadquarter": true
					}
				]
			}`,
		},
		{
			Name:        "Valid country code with single bank",
			CountryISO2: "lu",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCodes := []models.SwiftCode{
					{
						SwiftCode:     "BCEELULL",
						BankName:      "Banque et Caisse d'Epargne de l'Etat",
						CountryISO2:   "LU",
						CountryName:   "Luxembourg",
						Address:       "1, Place de Metz, Luxembourg",
						IsHeadquarter: true,
					},
				}
				repo.On("FindByCountryISO2", mock.Anything, "LU").Return(
					swiftCodes,
					"Luxembourg",
					nil,
				)
			},
			ExpectedStatus: http.StatusOK,
			ExpectedResponse: `{
				"countryISO2": "LU",
				"countryName": "Luxembourg",
				"swiftCodes": [
					{
						"swiftCode": "BCEELULL",
						"bankName": "Banque et Caisse d'Epargne de l'Etat",
						"countryISO2": "LU",
						"countryName": "Luxembourg",
						"address": "1, Place de Metz, Luxembourg",
						"isHeadquarter": true
					}
				]
			}`,
		},
		{
			Name:             "Invalid country code length",
			CountryISO2:      "USA",
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name:             "Invalid country code format (numbers)",
			CountryISO2:      "12",
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name:        "Country with no banks",
			CountryISO2: "ZZ",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("FindByCountryISO2", mock.Anything, "ZZ").Return(
					[]models.SwiftCode{},
					"",
					mongo.ErrNoDocuments,
				)
			},
			ExpectedStatus:   http.StatusNotFound,
			ExpectedResponse: `{"message":"No SWIFT codes found for this country"}`,
		},
	}
}
