package test_cases

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"swift-codes-api/models"
	mockRepo "swift-codes-api/repositories/mock"
)

type DeleteSwiftCodeTestCase struct {
	Name             string
	SwiftCode        string
	SetupMocks       func(repository *mockRepo.SwiftRepository)
	ExpectedStatus   int
	ExpectedResponse string
}

func GetDeleteSwiftCodeTestCases() []DeleteSwiftCodeTestCase {
	return []DeleteSwiftCodeTestCase{
		{
			Name:      "Successful deletion",
			SwiftCode: "ABCDUS12XXX",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCode := &models.SwiftCode{
					SwiftCode:     "ABCDUS12XXX",
					BankName:      "Bank of America",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 Main St, New York",
					IsHeadquarter: true,
				}
				repo.On("FindByCode", mock.Anything, "ABCDUS12XXX").Return(swiftCode, nil)
				repo.On("DeleteSwiftCode", mock.Anything, "ABCDUS12XXX").Return(nil)
			},
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"message":"SWIFT code deleted successfully"}`,
		},
		{
			Name:      "Invalid SWIFT code format",
			SwiftCode: "INVALID",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {

			},
			ExpectedStatus:   http.StatusBadRequest,
			ExpectedResponse: `{"message":"Invalid SWIFT code format"}`,
		},
		{
			Name:      "SWIFT code not found",
			SwiftCode: "ABCDJP12XXX",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "ABCDJP12XXX").Return(nil, mongo.ErrNoDocuments)
			},
			ExpectedStatus:   http.StatusNotFound,
			ExpectedResponse: `{"message":"SWIFT code not found"}`,
		},
		{
			Name:      "Delete operation failed",
			SwiftCode: "DEUTDE11XXX",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCode := &models.SwiftCode{
					SwiftCode:     "DEUTDE11XXX",
					BankName:      "Deutsche Bank",
					CountryISO2:   "DE",
					CountryName:   "Germany",
					Address:       "456 Main St, Berlin",
					IsHeadquarter: true,
					SwiftPrefix:   "DEUTDE",
				}
				repo.On("FindByCode", mock.Anything, "DEUTDE11XXX").Return(swiftCode, nil)
				repo.On("DeleteSwiftCode", mock.Anything, "DEUTDE11XXX").Return(errors.New("database error"))
			},
			ExpectedStatus:   http.StatusInternalServerError,
			ExpectedResponse: `{"message":"Failed to delete SWIFT code"}`,
		},
		{
			Name:      "Delete headquarter SWIFT code",
			SwiftCode: "BARCGB22XXX",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCode := &models.SwiftCode{
					SwiftCode:     "BARCGB22XXX",
					BankName:      "Barclays Bank",
					CountryISO2:   "GB",
					CountryName:   "United Kingdom",
					Address:       "1 Churchill Place, London",
					IsHeadquarter: true,
					SwiftPrefix:   "BARCGB",
				}
				repo.On("FindByCode", mock.Anything, "BARCGB22XXX").Return(swiftCode, nil)
				repo.On("DeleteSwiftCode", mock.Anything, "BARCGB22XXX").Return(nil)
			},
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"message":"SWIFT code deleted successfully"}`,
		},
	}
}
