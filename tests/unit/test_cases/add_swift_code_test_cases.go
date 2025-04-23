package test_cases

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"net/http"
	"swift-codes-api/models"
	mockRepo "swift-codes-api/repositories/mock"
)

type AddSwiftCodeTestCase struct {
	Name             string
	RequestBody      string
	SetupMocks       func(repository *mockRepo.SwiftRepository)
	ExpectedStatus   int
	ExpectedResponse string
}

func GetAddSwiftCodeTestCases() []AddSwiftCodeTestCase {
	return []AddSwiftCodeTestCase{
		{
			Name: "valid swift code",
			RequestBody: `{
				"swiftCode": "AAAAUS33",
				"bankName": "BANK OF AMERICA",
				"countryISO2": "us",
				"countryName": "United States",
				"address": "123 MAIN ST, NEW YORK",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.MatchedBy(func(sc models.SwiftCode) bool {
					return sc.SwiftCode == "AAAAUS33" && sc.CountryISO2 == "US"
				})).Return(nil)
			},
			ExpectedStatus:   http.StatusCreated,
			ExpectedResponse: `{"message":"SWIFT code added successfully"}`,
		},
		{
			Name: "missing required fields",
			RequestBody: `{
				"swiftCode": "AAAAUS33",
				"bankName": "",
				"countryISO2": "US",
				"countryName": "United States",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Missing required fields"}`,
		},
		{
			Name: "invalid country code format",
			RequestBody: `{
				"swiftCode": "AAAAUS33",
				"bankName": "BANK OF AMERICA",
				"countryISO2": "USA",
				"countryName": "United States",
				"address": "123 MAIN ST, NEW YORK",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name: "invalid swift code format - too short",
			RequestBody: `{
				"swiftCode": "AAAA",
				"bankName": "BANK OF AMERICA",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 MAIN ST, NEW YORK",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid SWIFT code format. Must be 8 or 11 characters"}`,
		},
		{
			Name: "swift code already exists",
			RequestBody: `{
				"swiftCode": "AAAAUS33",
				"bankName": "BANK OF AMERICA",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 MAIN ST, NEW YORK",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.MatchedBy(func(sc models.SwiftCode) bool {
					return sc.SwiftCode == "AAAAUS33"
				})).Return(errors.New("SWIFT code AAAAUS33 already exists"))
			},
			ExpectedStatus:   http.StatusConflict,
			ExpectedResponse: `{"message":"SWIFT code AAAAUS33 already exists"}`,
		},
		{
			Name: "server error",
			RequestBody: `{
				"swiftCode": "AAAAUS33",
				"bankName": "BANK OF AMERICA",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 MAIN ST, NEW YORK",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			ExpectedStatus:   http.StatusInternalServerError,
			ExpectedResponse: `{"message":"Failed to add SWIFT code"}`,
		},
	}
}
