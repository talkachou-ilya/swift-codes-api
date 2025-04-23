package test_cases

import (
	"errors"
	"net/http"

	"github.com/stretchr/testify/mock"

	"swift-codes-api/models"
	mockRep "swift-codes-api/repositories/mock"
)

type SwiftCodeTestCase struct {
	Name               string
	SwiftCode          string
	SetupMocks         func(repo *mockRep.SwiftRepository)
	ExpectedStatusCode int
	ExpectedResponse   string
}

func GetSwiftCodeTestCases() []SwiftCodeTestCase {
	return []SwiftCodeTestCase{
		{
			Name:      "Valid SWIFT code - non-headquarter",
			SwiftCode: "VALIDSWIFT",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "VALIDSWIFT").Return(&models.SwiftCode{
					SwiftCode:     "VALIDSWIFT",
					BankName:      "Test Bank",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 Test St",
					IsHeadquarter: false,
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"123 Test St","bankName":"Test Bank","countryISO2":"US","countryName":"United States","isHeadquarter":false,"swiftCode":"VALIDSWIFT"}`,
		},
		{
			Name:      "Valid SWIFT code - headquarter with branches",
			SwiftCode: "HQS1",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "HQS1").Return(&models.SwiftCode{
					SwiftCode:     "HQS1",
					BankName:      "HQ Bank",
					CountryISO2:   "FR",
					CountryName:   "France",
					Address:       "456 HQ St",
					IsHeadquarter: true,
					SwiftPrefix:   "HQS",
				}, nil)
				repo.On("FindBranchesByPrefix", mock.Anything, "HQS").Return([]models.SwiftCode{
					{SwiftCode: "HQS2", BankName: "Branch A"},
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"456 HQ St","bankName":"HQ Bank","branches":[{"address":"","bankName":"Branch A","countryISO2":"","countryName":"","isHeadquarter":false,"swiftCode":"HQS2"}],"countryISO2":"FR","countryName":"France","isHeadquarter":true,"swiftCode":"HQS1"}`,
		},
		{
			Name:      "Invalid SWIFT code",
			SwiftCode: "INVALIDSWIFT",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "INVALIDSWIFT").Return(nil, errors.New("not found"))
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"SWIFT code not found"}`,
		},
		{
			Name:      "Headquarter without branches",
			SwiftCode: "HQNO",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "HQNO").Return(&models.SwiftCode{
					SwiftCode:     "HQNO",
					BankName:      "HQ No Branch Bank",
					CountryISO2:   "DE",
					CountryName:   "Germany",
					Address:       "789 Main St",
					IsHeadquarter: true,
					SwiftPrefix:   "HQN",
				}, nil)
				repo.On("FindBranchesByPrefix", mock.Anything, "HQN").Return(nil, errors.New("no branches found"))
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"789 Main St","bankName":"HQ No Branch Bank","countryISO2":"DE","countryName":"Germany","isHeadquarter":true,"swiftCode":"HQNO"}`,
		},
	}
}
