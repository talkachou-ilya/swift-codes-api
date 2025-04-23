package unit

import (
	"net/http"
	"net/http/httptest"
	"swift-codes-api/tests/unit/test_cases"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"swift-codes-api/handlers"
	"swift-codes-api/internal/config"
	mockRepos "swift-codes-api/repositories/mock"
)

func TestGetSwiftCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tc := range test_cases.GetSwiftCodeTestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo := new(mockRepos.SwiftRepository)
			tc.SetupMocks(mockRepo)

			handler := handlers.NewSwiftHandler(config.Config{}, mockRepo)
			router := setupRouter(handler)

			response := performRequest(router, tc.SwiftCode)

			assert.Equal(t, tc.ExpectedStatusCode, response.Code)
			assert.JSONEq(t, tc.ExpectedResponse, response.Body.String())

			mockRepo.AssertExpectations(t)
		})
	}
}

func setupRouter(handler *handlers.SwiftCodesHandler) *gin.Engine {
	router := gin.Default()
	router.GET("/swift/:swift-code", handler.GetSwiftCode)
	return router
}

func performRequest(router *gin.Engine, swiftCode string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/swift/"+swiftCode, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
