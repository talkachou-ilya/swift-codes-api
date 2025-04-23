package unit

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"swift-codes-api/handlers"
	"swift-codes-api/internal/config"
	mockRepos "swift-codes-api/repositories/mock"
	"swift-codes-api/tests/unit/test_cases"
	"testing"
)

func TestGetSwiftCodesByCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tc := range test_cases.GetCountryTestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo := new(mockRepos.SwiftRepository)
			tc.SetupMocks(mockRepo)

			handler := handlers.NewSwiftHandler(config.Config{}, mockRepo)

			router := gin.Default()
			router.GET("/swift-codes/country/:countryISO2code", handler.GetSwiftCodesByCountry)

			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/swift-codes/country/"+tc.CountryISO2, nil)
			router.ServeHTTP(recorder, req)

			respBody := recorder.Body.String()
			require.Equal(t, tc.ExpectedStatus, recorder.Code)
			require.JSONEq(t, tc.ExpectedResponse, respBody)
			mockRepo.AssertExpectations(t)
		})
	}
}
