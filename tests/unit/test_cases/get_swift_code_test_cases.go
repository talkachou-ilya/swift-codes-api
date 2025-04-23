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
			SwiftCode: "ABCDUS12XXX",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "ABCDUS12XXX").Return(&models.SwiftCode{
					SwiftCode:     "ABCDUS12XXX",
					BankName:      "Test Bank",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 Test St",
					IsHeadquarter: false,
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"123 Test St","bankName":"Test Bank","countryISO2":"US","countryName":"United States","isHeadquarter":false,"swiftCode":"ABCDUS12XXX"}`,
		},
		{
			Name:      "Valid SWIFT code - headquarter with branches",
			SwiftCode: "DEUTDE11XXX",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "DEUTDE11XXX").Return(&models.SwiftCode{
					SwiftCode:     "DEUTDE11XXX",
					BankName:      "Deutsche Bank",
					CountryISO2:   "DE",
					CountryName:   "Germany",
					Address:       "456 Main St, Berlin",
					IsHeadquarter: true,
					SwiftPrefix:   "DEUTDE",
				}, nil)
				repo.On("FindBranchesByPrefix", mock.Anything, "DEUTDE").Return([]models.SwiftCode{
					{SwiftCode: "DEUTDE22XXX", BankName: "Deutsche Bank Branch", CountryISO2: "DE", CountryName: "Germany", Address: "789 Branch St, Munich"},
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"456 Main St, Berlin","bankName":"Deutsche Bank","branches":[{"address":"789 Branch St, Munich","bankName":"Deutsche Bank Branch","countryISO2":"DE","countryName":"Germany","isHeadquarter":false,"swiftCode":"DEUTDE22XXX"}],"countryISO2":"DE","countryName":"Germany","isHeadquarter":true,"swiftCode":"DEUTDE11XXX"}`,
		},
		{
			Name:      "Valid SWIFT code - not found",
			SwiftCode: "ABCDEF12XXX",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "ABCDEF12XXX").Return(nil, errors.New("not found"))
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   `{"message":"SWIFT code not found"}`,
		},
		{
			Name:      "Invalid SWIFT code format",
			SwiftCode: "INVALID",
			SetupMocks: func(repo *mockRep.SwiftRepository) {

			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   `{"message":"Invalid SWIFT code format."}`,
		},
		{
			Name:      "Headquarter without branches",
			SwiftCode: "BARCGB22XXX",
			SetupMocks: func(repo *mockRep.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "BARCGB22XXX").Return(&models.SwiftCode{
					SwiftCode:     "BARCGB22XXX",
					BankName:      "Barclays Bank",
					CountryISO2:   "GB",
					CountryName:   "United Kingdom",
					Address:       "1 Churchill Place, London",
					IsHeadquarter: true,
					SwiftPrefix:   "BARCGB",
				}, nil)
				repo.On("FindBranchesByPrefix", mock.Anything, "BARCGB").Return(nil, errors.New("no branches found"))
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   `{"address":"1 Churchill Place, London","bankName":"Barclays Bank","countryISO2":"GB","countryName":"United Kingdom","isHeadquarter":true,"swiftCode":"BARCGB22XXX"}`,
		},
	}
}
