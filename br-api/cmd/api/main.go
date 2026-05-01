package main

import (
	"os"

	"github.com/joho/godotenv"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.New(database)
	logger.Log.Info("starting server", "port", port)
	if err := r.Run(":" + port); err != nil {
		logger.Log.Error("server failed", "error", err)
		os.Exit(1)
	}
}
