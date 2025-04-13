package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/infrastructure/config"
	"github.com/merdernoty/anime-service/internal/infrastructure/database"
	"github.com/merdernoty/anime-service/internal/infrastructure/log"
	httpServer "github.com/merdernoty/anime-service/internal/interfaces/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	logger := log.NewLogger(log.Config{
		Format:  cfg.Logger.Format,
		Level:   cfg.Logger.Level,
		Nocolor: cfg.App.Environment != "development",
	})

	log.SetStandartLogger(logger)

	// Подключение к базе данных
	db, err := database.NewConnector(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Error closing database connection", map[string]interface{}{"error": err.Error()})
		}
	}()

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// userRepo := postgres.NewUserRepository(db)
	// animeRepo := postgres.NewAnimeRepository(db)

	// userService := services.NewUserService(userRepo, logger)
	// animeService := services.NewAnimeService(animeRepo, logger)

	// authMiddleware := middleware.NewAuthMiddleware(cfg.Auth)
	// loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	// userController := controllers.NewUserController(userService)
	// animeController := controllers.NewAnimeController(animeService)

	server := httpServer.NewServer(
		cfg,
		// userController,
		// animeController,
		// authMiddleware,
		// loggerMiddleware,
	)

	logger.Info("Starting server", map[string]interface{}{
		"port": cfg.HTTP.Port,
		"env":  cfg.App.Environment,
	})

	if err := server.Start(); err != nil {
		logger.Error("Server error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
