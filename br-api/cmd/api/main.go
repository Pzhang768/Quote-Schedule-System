package main

import (
	"log"
	"os"
	// "github.com/melfish/br-api/internal/db"

	"github.com/joho/godotenv"
	"github.com/melfish/br-api/internal/router"
)

func main() {
	_ = godotenv.Load(".env.local")

	// database, err := db.Connect(dsn)
	// if err != nil {
	// 	log.Fatalf("failed to connect to database: %v", err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.New(nil)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
