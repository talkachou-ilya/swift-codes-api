package unit

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"swift-codes-api/handlers"
	"swift-codes-api/internal/config"
	mockRepos "swift-codes-api/repositories/mock"
	"swift-codes-api/tests/unit/test_cases"
	"testing"
)

func TestAddSwiftCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tc := range test_cases.GetAddSwiftCodeTestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			mockRepo := new(mockRepos.SwiftRepository)
			tc.SetupMocks(mockRepo)

			handler := handlers.NewSwiftHandler(config.Config{}, mockRepo)
			router := gin.Default()
			router.POST("/swift-codes", handler.AddSwiftCode)

			req := httptest.NewRequest(http.MethodPost, "/swift-codes", bytes.NewBufferString(tc.RequestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.ExpectedStatus, w.Code)
			assert.JSONEq(t, tc.ExpectedResponse, w.Body.String())

			mockRepo.AssertExpectations(t)
		})
	}
}
