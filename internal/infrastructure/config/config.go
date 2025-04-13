// internal/infrastructure/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/merdernoty/anime-service/internal/infrastructure/database"
)

type Config struct {
	App        AppConfig        // Общие настройки приложения
	HTTP       HTTPConfig       // Настройки HTTP-сервера
	Database   database.Config  // Настройки базы данных (из пакета database)
	Logger     LoggerConfig     // Настройки логирования
	Auth       AuthConfig       // Настройки аутентификации
	Cache      CacheConfig      // Настройки кэширования
	Pagination PaginationConfig // Настройки пагинации
}

type AppConfig struct {
	Name        string // Название приложения
	Environment string // Окружение (development, staging, production)
	Debug       bool   // Режим отладки
	TimeZone    string // Временная зона
}

type HTTPConfig struct {
	Host               string        // Хост для привязки сервера
	Port               int           // Порт для привязки сервера
	ReadTimeout        time.Duration // Таймаут для чтения запроса
	WriteTimeout       time.Duration // Таймаут для записи ответа
	IdleTimeout        time.Duration // Таймаут для неактивных соединений
	MaxHeaderBytes     int           // Максимальный размер заголовков
	TrustedProxies     []string      // Доверенные прокси для корректного получения IP
	AllowedCORSOrigins []string      // Разрешенные источники для CORS
}

type LoggerConfig struct {
	Level        string // Уровень логирования (debug, info, warn, error)
	Format       string // Формат логов (json, text)
	Output       string // Вывод логов (stdout, file)
	FilePath     string // Путь к файлу логов
	MaxSize      int    // Максимальный размер файла логов в МБ
	MaxBackups   int    // Максимальное количество файлов логов
	MaxAge       int    // Максимальное время хранения файлов логов в днях
	Compress     bool   // Сжимать ли файлы логов
	ReportCaller bool   // Логировать ли вызывающий код
}

type AuthConfig struct {
	JWTSecret              string        // Секретный ключ для JWT
	AccessTokenExpiration  time.Duration // Время жизни access token
	RefreshTokenExpiration time.Duration // Время жизни refresh token
	TokenIssuer            string        // Издатель токена
}

type CacheConfig struct {
	Driver            string        // Драйвер кэша (redis, memory)
	Host              string        // Хост Redis
	Port              int           // Порт Redis
	Password          string        // Пароль Redis
	DB                int           // Номер базы данных Redis
	DefaultExpiration time.Duration // Время жизни кэша по умолчанию
}

type PaginationConfig struct {
	DefaultLimit int // Лимит по умолчанию
	MaxLimit     int // Максимальный лимит
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbParams := make(map[string]string)
	if sslMode := getEnv("DB_SSLMODE", "disable"); sslMode != "" {
		dbParams["sslmode"] = sslMode
	}
	if appName := getEnv("APP_NAME", "anime-service"); appName != "" {
		dbParams["application_name"] = appName
	}

	config := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "anime-service"),
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
			TimeZone:    getEnv("APP_TIMEZONE", "UTC"),
		},
		HTTP: HTTPConfig{
			Host:               getEnv("HTTP_HOST", ""),
			Port:               getEnvAsInt("HTTP_PORT", 8080),
			ReadTimeout:        getEnvAsDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:       getEnvAsDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:        getEnvAsDuration("HTTP_IDLE_TIMEOUT", 120*time.Second),
			MaxHeaderBytes:     getEnvAsInt("HTTP_MAX_HEADER_BYTES", 1<<20), // 1 MB
			TrustedProxies:     getEnvAsSlice("HTTP_TRUSTED_PROXIES", []string{}),
			AllowedCORSOrigins: getEnvAsSlice("HTTP_ALLOWED_CORS_ORIGINS", []string{"*"}),
		},
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
			Params:   dbParams,
		},
		Logger: LoggerConfig{
			Level:        getEnv("LOG_LEVEL", "info"),
			Format:       getEnv("LOG_FORMAT", "json"),
			Output:       getEnv("LOG_OUTPUT", "stdout"),
			FilePath:     getEnv("LOG_FILE_PATH", "logs/anime-service.log"),
			MaxSize:      getEnvAsInt("LOG_MAX_SIZE", 10),
			MaxBackups:   getEnvAsInt("LOG_MAX_BACKUPS", 3),
			MaxAge:       getEnvAsInt("LOG_MAX_AGE", 7),
			Compress:     getEnvAsBool("LOG_COMPRESS", true),
			ReportCaller: getEnvAsBool("LOG_REPORT_CALLER", false),
		},
		Auth: AuthConfig{
			JWTSecret:              getEnv("AUTH_JWT_SECRET", "your-secret-key"),
			AccessTokenExpiration:  getEnvAsDuration("AUTH_ACCESS_TOKEN_EXPIRATION", 15*time.Minute),
			RefreshTokenExpiration: getEnvAsDuration("AUTH_REFRESH_TOKEN_EXPIRATION", 7*24*time.Hour),
			TokenIssuer:            getEnv("AUTH_TOKEN_ISSUER", "anime-service"),
		},
		Cache: CacheConfig{
			Driver:            getEnv("CACHE_DRIVER", "memory"),
			Host:              getEnv("CACHE_REDIS_HOST", "localhost"),
			Port:              getEnvAsInt("CACHE_REDIS_PORT", 6379),
			Password:          getEnv("CACHE_REDIS_PASSWORD", ""),
			DB:                getEnvAsInt("CACHE_REDIS_DB", 0),
			DefaultExpiration: getEnvAsDuration("CACHE_DEFAULT_EXPIRATION", 5*time.Minute),
		},
		Pagination: PaginationConfig{
			DefaultLimit: getEnvAsInt("PAGINATION_DEFAULT_LIMIT", 10),
			MaxLimit:     getEnvAsInt("PAGINATION_MAX_LIMIT", 100),
		},
	}

	if err := config.Database.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
