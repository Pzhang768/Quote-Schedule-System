// @title           Quote Scheduler API
// @version         1.0
// @description     Quote scheduling and notification system
// @BasePath        /api/v1

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/melfish/br-api/docs"
	"github.com/melfish/br-api/internal/db"
	"github.com/melfish/br-api/internal/logger"
	"github.com/melfish/br-api/internal/router"
)

func main() {
	_ = godotenv.Load(".env.local")
	logger.Init()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:root@tcp(localhost:3307)/brix?parseTime=true"
	}

	database, err := db.Connect(dsn)
	if err != nil {
		logger.Log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	logger.Log.Info("database connected")

	if err := db.Seed(database); err != nil {
		logger.Log.Error("seed failed", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.SwaggerInfo.Host = "localhost:" + port

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000"
	}

	r := router.New(database, corsOrigin)
	logger.Log.Info("starting server", "port", port)
	if err := r.Run(":" + port); err != nil {
		logger.Log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
