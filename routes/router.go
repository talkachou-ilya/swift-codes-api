package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"swift-codes-api/handlers"
	"swift-codes-api/internal/config"
	mongo2 "swift-codes-api/repositories/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database, cfg config.Config) {
	h := handlers.NewSwiftHandler(cfg, mongo2.NewSwiftRepository(db))

	v1 := r.Group("/v1/swift-codes")
	{
		v1.GET("/:swift-code", h.GetSwiftCode)
	}
}
