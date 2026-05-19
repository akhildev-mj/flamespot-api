package service

import (
	"flamespot-api/config"
	"flamespot-api/model"
	"flamespot-api/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MenuService interface {
	GetMenu() ([]model.MenuItem, error)
}

type menuService struct {
	db *mongo.Database
}

func NewMenuService(db *mongo.Database) MenuService {
	return &menuService{db: db}
}

func (s *menuService) GetMenu() ([]model.MenuItem, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionMenu)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menuItems []model.MenuItem
	if err := cursor.All(ctx, &menuItems); err != nil {
		return nil, err
	}

	if menuItems == nil {
		menuItems = []model.MenuItem{}
	}

	return menuItems, nil
}
