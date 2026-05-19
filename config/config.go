package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	// Database
	CollectionMenu   = "menu"
	CollectionOrders = "orders"
	DatabaseName     = "flamespot_db"

	// Caching & Pagination
	DefaultTimeout  = 10 * time.Second
	CacheExpiration = 30 * 24 * time.Hour
	LimitDefault    = 10

	// Server Limits
	MaxBodySize  = 1 * 1024 * 1024
	IdleTimeout  = 10 * time.Second
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second

	// CORS Constants (Stored safely as constant strings)
	CORSMethods = "GET,POST,OPTIONS"
	CORSHeaders = "Origin,Content-Type,Accept,Authorization"
)

type Config struct {
	MongoURI           string
	Port               string
	JWTSecret          string
	JWTExpiration      time.Duration
	CORSAllowedOrigins []string
}

func LoadConfig() Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	jwtExpiry := os.Getenv("JWT_EXPIRATION_DAYS")
	if jwtExpiry == "" {
		log.Fatal("JWT_EXPIRATION_DAYS environment variable is not set")
	}

	expDays, err := strconv.Atoi(jwtExpiry)
	if err != nil || expDays <= 0 {
		log.Fatalf("JWT_EXPIRATION_DAYS must be a valid positive number, got: '%s'", jwtExpiry)
	}

	var allowedOrigins []string
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")

	if originsStr == "" {
		log.Println("Warning: CORS_ALLOWED_ORIGINS is not set. Defaulting to empty list.")
		allowedOrigins = []string{}
	} else {
		rawOrigins := strings.Split(originsStr, ",")
		for _, o := range rawOrigins {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(o))
		}
	}

	return Config{
		Port:               port,
		MongoURI:           mongoURI,
		JWTSecret:          jwtSecret,
		JWTExpiration:      time.Duration(expDays) * 24 * time.Hour,
		CORSAllowedOrigins: allowedOrigins,
	}
}
