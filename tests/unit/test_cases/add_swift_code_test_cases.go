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
			Name: "Valid SWIFT code",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
				"bankName": "Bank of America",
				"countryISO2": "us",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.MatchedBy(func(sc models.SwiftCode) bool {
					return sc.SwiftCode == "ABCDUS12XXX" && sc.CountryISO2 == "US"
				})).Return(nil)
			},
			ExpectedStatus:   http.StatusCreated,
			ExpectedResponse: `{"message":"SWIFT code added successfully"}`,
		},
		{
			Name: "Missing required fields",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
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
			Name: "Invalid country code format",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
				"bankName": "Bank of America",
				"countryISO2": "USA",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid country code format. Must be a 2-letter ISO country code"}`,
		},
		{
			Name: "Invalid SWIFT code format",
			RequestBody: `{
				"swiftCode": "INVALID",
				"bankName": "Bank of America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid SWIFT code format. Must be 11 characters and follow proper format"}`,
		},
		{
			Name: "Country code mismatch in SWIFT code",
			RequestBody: `{
				"swiftCode": "ABCDJP12XXX",
				"bankName": "Bank of America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Country code in SWIFT code does not match the provided country code"}`,
		},
		{
			Name: "SWIFT code already exists",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
				"bankName": "Bank of America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.MatchedBy(func(sc models.SwiftCode) bool {
					return sc.SwiftCode == "ABCDUS12XXX"
				})).Return(errors.New("SWIFT code ABCDUS12XXX already exists"))
			},
			ExpectedStatus:   http.StatusConflict,
			ExpectedResponse: `{"message":"SWIFT code ABCDUS12XXX already exists"}`,
		},
		{
			Name: "Server error",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
				"bankName": "Bank of America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true
			}`,
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("AddSwiftCode", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			ExpectedStatus:   http.StatusInternalServerError,
			ExpectedResponse: `{"message":"Failed to add SWIFT code"}`,
		},
		{
			Name: "Invalid JSON format",
			RequestBody: `{
				"swiftCode": "ABCDUS12XXX",
				"bankName": "Bank of America",
				"countryISO2": "US",
				"countryName": "United States",
				"address": "123 Main St, New York",
				"isHeadquarter": true,
			}`,
			SetupMocks:       func(repo *mockRepo.SwiftRepository) {},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid request format"}`,
		},
	}
}
