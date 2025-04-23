package integration

import (
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

func TestGetSwiftCode_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Load()
	testApp := app.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pingErr := testApp.MongoDB.Client().Ping(ctx, nil)
	log.Printf("MongoDB connection test: %v", pingErr)

	for _, tc := range test_cases.GetIntegrationTestCases() {
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
			log.Println(documents)

			response := performRequest(testApp.Router, tc.SwiftCode)

			assert.Equal(t, tc.ExpectedStatusCode, response.Code)
			assert.JSONEq(t, tc.ExpectedResponse, response.Body.String())

			if err := testApp.MongoDB.Collection("swift-codes").Drop(testCtx); err != nil {
				t.Logf("Warning: Failed to drop collection: %v", err)
			}
		})
	}

	if err := testApp.Mongo.Disconnect(ctx); err != nil {
		t.Logf("Warning: Failed to disconnect MongoDB client: %v", err)
	}
}

func performRequest(router *gin.Engine, swiftCode string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/v1/swift-codes/"+swiftCode, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
