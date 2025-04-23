package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"swift-codes-api/handlers"
	"swift-codes-api/internal/config"
	repos "swift-codes-api/repositories/mongo"
)

func SetupRoutes(r *gin.Engine, db *mongo.Database, cfg config.Config) {
	h := handlers.NewSwiftHandler(cfg, repos.NewSwiftRepository(db))

	v1 := r.Group("/v1/swift-codes")
	{
		v1.GET("/:swift-code", h.GetSwiftCode)
		v1.GET("/country/:countryISO2code", h.GetSwiftCodesByCountry)
		v1.POST("", h.AddSwiftCode)
		v1.DELETE("/:swift-code", h.DeleteSwiftCode)
	}
}
