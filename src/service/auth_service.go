package service

import (
	"errors"
	"time"

	"flamespot-api/src/config"
	"flamespot-api/src/model"
	"flamespot-api/src/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (string, model.User, error)
}

type authService struct {
	db  *mongo.Database
	cfg config.Config
}

func NewAuthService(db *mongo.Database, cfg config.Config) AuthService {
	return &authService{db: db, cfg: cfg}
}

func (s *authService) Login(username, password string) (string, model.User, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	var user model.User
	err := s.db.Collection(config.CollectionUsers).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", model.User{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", model.User{}, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user": user.Username,
		"exp":  time.Now().Add(s.cfg.JWTExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWTSecret))

	return signedToken, user, err
}
