package http
import (	

	"net/http"
    "github.com/gin-gonic/gin"
    
    "anime-service/internal/interfaces/http/controllers"
    "anime-service/internal/interfaces/http/middleware"
    "anime-service/internal/interfaces/http/routes"
    "anime-service/internal/infrastructure/config"
)
type Server struct {
    router     *gin.Engine
    config     *config.Config
    httpServer *http.Server
}