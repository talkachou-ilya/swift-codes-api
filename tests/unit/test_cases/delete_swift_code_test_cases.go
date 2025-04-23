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
			Name:      "successful deletion",
			SwiftCode: "AAAAUS33",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCode := &models.SwiftCode{
					SwiftCode:     "AAAAUS33",
					BankName:      "BANK OF AMERICA",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 MAIN ST, NEW YORK",
					IsHeadquarter: true,
				}
				repo.On("FindByCode", mock.Anything, "AAAAUS33").Return(swiftCode, nil)
				repo.On("DeleteSwiftCode", mock.Anything, "AAAAUS33").Return(nil)
			},
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"message":"SWIFT code deleted successfully"}`,
		},
		{
			Name:      "swift code not found",
			SwiftCode: "NOTFOUND",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				repo.On("FindByCode", mock.Anything, "NOTFOUND").Return(nil, mongo.ErrNoDocuments)
			},
			ExpectedStatus:   http.StatusNotFound,
			ExpectedResponse: `{"message":"SWIFT code not found"}`,
		},
		{
			Name:      "delete operation failed",
			SwiftCode: "AAAAUS33",
			SetupMocks: func(repo *mockRepo.SwiftRepository) {
				swiftCode := &models.SwiftCode{
					SwiftCode:     "AAAAUS33",
					BankName:      "BANK OF AMERICA",
					CountryISO2:   "US",
					CountryName:   "United States",
					Address:       "123 MAIN ST, NEW YORK",
					IsHeadquarter: true,
				}
				repo.On("FindByCode", mock.Anything, "AAAAUS33").Return(swiftCode, nil)
				repo.On("DeleteSwiftCode", mock.Anything, "AAAAUS33").Return(errors.New("database error"))
			},
			ExpectedStatus:   http.StatusInternalServerError,
			ExpectedResponse: `{"message":"Failed to delete SWIFT code"}`,
		},
	}
}
