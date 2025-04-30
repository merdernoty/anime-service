package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merdernoty/anime-service/internal/infrastructure/config"
	"github.com/merdernoty/anime-service/internal/interfaces/http/controllers"
	"github.com/merdernoty/anime-service/internal/interfaces/http/middleware"
	"github.com/merdernoty/anime-service/internal/interfaces/http/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/merdernoty/anime-service/docs"
)
type Server struct {
    router     *gin.Engine
    config     *config.Config
    httpServer *http.Server
}

func NewServer(
    config *config.Config,
    authConttroler *controllers.AuthController,
    animeController *controllers.AnimeController,
    authMiddleware middleware.AuthMiddleware,
) *Server {
    router := gin.New()
    if config.App.Environment == "development" {
        router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }
    router.Use(gin.Recovery())
    service := &routes.Service{
        AuthController: authConttroler,
        AnimeController: animeController,
    }
    routes.SetupRoutes(router, service, &authMiddleware)
    
    // Создание HTTP-сервера
    httpServer := &http.Server{
        Addr:         fmt.Sprintf(":%d", config.HTTP.Port),
        Handler:      router,
        ReadTimeout:  time.Duration(config.HTTP.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(config.HTTP.WriteTimeout) * time.Second,
        IdleTimeout:  time.Duration(config.HTTP.IdleTimeout) * time.Second,
    }
    
    return &Server{
        router:     router,
        config:     config,
        httpServer: httpServer,
    }
}

func (s *Server) Start() error {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        log.Printf("HTTP server listening on %s", s.httpServer.Addr)
        if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    
    <-quit
    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    
    if err := s.httpServer.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited properly")
    return nil
}