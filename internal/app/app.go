package app

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"swift-codes-api/internal/config"
	"swift-codes-api/internal/db"
	"swift-codes-api/routes"
)

type App struct {
	Config  config.Config
	Router  *gin.Engine
	Mongo   *mongo.Client
	MongoDB *mongo.Database
}

func New(cfg config.Config) *App {
	dbClient := db.Connect(cfg.MongoURI)
	database := dbClient.Database(cfg.MongoDB)

	r := gin.Default()
	routes.SetupRoutes(r, database, cfg)

	return &App{
		Config:  cfg,
		Router:  r,
		Mongo:   dbClient,
		MongoDB: database,
	}
}

func Start(a *App) {
	err := a.Router.Run(":" + a.Config.Port)
	if err != nil {
		log.Fatal(err)
	}
}
