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
	CollectionUsers      = "users"
	CollectionCategories = "categories"
	CollectionMenu       = "menu"
	CollectionOrders     = "orders"

	DefaultTimeout  = 10 * time.Second
	CacheExpiration = 30 * 24 * time.Hour
	LimitDefault    = 10

	MaxBodySize  = 1 * 1024 * 1024
	IdleTimeout  = 10 * time.Second
	ReadTimeout  = 5 * time.Second
	WriteTimeout = 5 * time.Second

	CORSMethods = "GET,POST,PATCH,OPTIONS"
	CORSHeaders = "Origin,Content-Type,Accept,Authorization"
)

type Config struct {
	MongoURI           string
	DatabaseName       string
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

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is not set")
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
		log.Fatalf("JWT_EXPIRATION_DAYS must be a positive number, got: '%s'", jwtExpiry)
	}

	var allowedOrigins []string
	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsStr != "" {
		rawOrigins := strings.Split(originsStr, ",")
		for _, o := range rawOrigins {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(o))
		}
	}

	return Config{
		Port:               port,
		MongoURI:           mongoURI,
		DatabaseName:       dbName,
		JWTSecret:          jwtSecret,
		JWTExpiration:      time.Duration(expDays) * 24 * time.Hour,
		CORSAllowedOrigins: allowedOrigins,
	}
}
