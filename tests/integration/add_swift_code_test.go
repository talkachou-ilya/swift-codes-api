package integration

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"net/http/httptest"
	"swift-codes-api/internal/app"
	"swift-codes-api/internal/config"
	"swift-codes-api/tests/integration/test_cases"
	"testing"
	"time"
)

func TestAddSwiftCode_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	testApp := app.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pingErr := testApp.MongoDB.Client().Ping(ctx, nil)
	log.Printf("MongoDB connection test: %v", pingErr)

	for _, tc := range test_cases.GetAddSwiftCodeIntegrationTestCases() {
		t.Run(tc.Name, func(t *testing.T) {
			testCtx, testCancel := context.WithTimeout(ctx, 5*time.Second)
			defer testCancel()

			if err := testApp.MongoDB.Collection("swift-codes").Drop(testCtx); err != nil {
				t.Logf("Warning: Failed to drop collection: %v", err)
			}

			tc.SetupData(testCtx, testApp.MongoDB)

			documents, err := testApp.MongoDB.Collection("swift-codes").CountDocuments(testCtx, bson.M{})
			if err != nil {
				t.Logf("Warning: Failed to count documents: %v", err)
			}
			log.Println("Initial document count:", documents)

			response := performAddSwiftCodeRequest(testApp.Router, tc.RequestBody)

			assert.Equal(t, tc.ExpectedStatusCode, response.Code)
			assert.JSONEq(t, tc.ExpectedResponse, response.Body.String())

			if !tc.CheckData(testCtx, testApp.MongoDB) {
				t.Error("Database state not as expected after request")
			}

			if err := testApp.MongoDB.Collection("swift-codes").Drop(testCtx); err != nil {
				t.Logf("Warning: Failed to drop collection: %v", err)
			}
		})
	}

	if err := testApp.Mongo.Disconnect(ctx); err != nil {
		t.Logf("Warning: Failed to disconnect MongoDB client: %v", err)
	}
}

func performAddSwiftCodeRequest(router *gin.Engine, requestBody string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/v1/swift-codes", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
