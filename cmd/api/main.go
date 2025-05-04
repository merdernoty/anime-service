package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/application/services"
	"github.com/merdernoty/anime-service/internal/domain/models"
	"github.com/merdernoty/anime-service/internal/infrastructure/api"
	"github.com/merdernoty/anime-service/internal/infrastructure/config"
	"github.com/merdernoty/anime-service/internal/infrastructure/database"
	"github.com/merdernoty/anime-service/internal/infrastructure/log"
	"github.com/merdernoty/anime-service/internal/infrastructure/repositories"
	httpServer "github.com/merdernoty/anime-service/internal/interfaces/http"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
	"github.com/merdernoty/anime-service/internal/interfaces/http/middleware"
	"github.com/merdernoty/anime-service/internal/interfaces/http/routes"
	"github.com/merdernoty/anime-service/pkg/auth"
)

//	@title						Anime Service API
//	@version					1.0
//	@description				API для сервиса аниме и управления пользовательскими списками
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Введите токен в формате: Bearer {token}
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
		if err := database.CloseDB(db); err != nil {
			logger.Error("Error closing database connection", map[string]interface{}{"error": err.Error()})
		}
	}()
    
	// Автомиграция моделей
	if err := database.AutoMigrate(db, 
        &models.User{},
        &models.UserAnime{},
    ); err != nil {
        logger.Error("Failed to auto-migrate database schema", map[string]interface{}{"error": err.Error()})
        os.Exit(1)
    }

	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	userRepo := repositories.NewUserRepository(db)
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get *sql.DB from *gorm.DB", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	userAnimeRepo := repositories.NewUserAnimeRepository(sqlDB, logger)

	jikanClient := api.NewJikanClient(logger)

	tokenMaker := auth.NewJWTTokenMaker(
		cfg.Auth.SecretKey, 
		cfg.Auth.TokenDuration,
	)

	authService := services.NewAuthService(
		userRepo, 
		logger,
		tokenMaker,
	)

	animeService := services.NewAnimeService(
		jikanClient,
		userAnimeRepo,
		logger,
	)
	
	authMiddleware := middleware.NewAuthMiddleware(tokenMaker, userRepo)
	authController := controllers.NewAuthController(authService)
	animeController := controllers.NewAnimeController(*animeService, logger)

	router := gin.Default()
	
	
	routes.SetupRoutes(
		router, 
		routes.NewService(
			authController,
			animeController,
		),
		authMiddleware,
	)

	server := httpServer.NewServer(
		cfg,
		authController,
		animeController,
		*authMiddleware,
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