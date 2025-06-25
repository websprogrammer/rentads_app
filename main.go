package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"rentads_app/database"
	"rentads_app/handlers"
	"time"
)

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to get current directory")
	}
	// Open or create the log file with required permissions
	file, err := os.OpenFile(
		filepath.Join(currentDir, "rentads_app.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Set the output destination for the standard logger
	log.SetOutput(file)
}

func main() {
	setUpLogging()

	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig, err := database.LoadConfig(".env")
	if err != nil {
		slog.Error("Error loading config", slog.String("error=", err.Error()))
	}
	dbConfig.MaxConns = 10
	dbConfig.MinConns = 2                        // Minimum connections in the pool, default is 0
	dbConfig.MaxConnLifeTime = 30 * time.Minute  // Maximum connection lifetime, default is one hour
	dbConfig.MaxConnIdleTime = 10 * time.Minute  // Maximum idle time, default is 30 minutes
	dbConfig.HealthCheckPeriod = 2 * time.Minute // Health check frequency, default is 60 seconds

	db, err := database.NewPg(rootCtx, dbConfig, database.WithPgxConfig(dbConfig))
	if err != nil {
		slog.Error("Error connecting to database", slog.String("error", err.Error()))
		panic(err)
	}
	defer db.Close()

	dbClient := &database.DB{
		Client: db,
	}
	dbClient.MonitorPoolStats()

	app := fiber.New()
	app.Get("/", handlers.GetAdverts(dbClient))
	app.Post("/add_notification", handlers.AddNotification(dbClient))
	log.Fatal(app.Listen(":3000"))

}
