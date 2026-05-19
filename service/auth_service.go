package service

import (
	"errors"
	"time"

	"flamespot-api/config"
	"flamespot-api/model"
	"flamespot-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, error)
}

type authService struct {
	db  *mongo.Database
	cfg config.Config
}

func NewAuthService(db *mongo.Database, cfg config.Config) AuthService {
	return &authService{db: db, cfg: cfg}
}

func (s *authService) Login(username, password string) (string, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	var user model.User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user": user.Username,
		"role": user.Role,
		"exp":  time.Now().Add(s.cfg.JWTExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
