package service

import (
	"errors"
	"time"

	"flamespot-api/src/config"
	"flamespot-api/src/model"
	"flamespot-api/src/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MenuService interface {
	GetMenu(isActiveFilter *bool, categoryIdFilter string, sortFilter string) ([]model.MenuItem, error)
	CreateMenuItem(item model.MenuItem) (model.MenuItem, error)
	UpdateMenuItem(id string, updates map[string]interface{}) (model.MenuItem, error)
}

type menuService struct {
	db *mongo.Database
}

func NewMenuService(db *mongo.Database) MenuService {
	return &menuService{db: db}
}

func (s *menuService) GetMenu(isActiveFilter *bool, categoryIdFilter string, sortFilter string) ([]model.MenuItem, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionMenu)

	filter := bson.M{}

	if isActiveFilter != nil {
		filter["isActive"] = *isActiveFilter
	}

	if categoryIdFilter != "" {
		if catID, err := bson.ObjectIDFromHex(categoryIdFilter); err == nil {
			filter["categoryId"] = catID
		}
	}

	opts := options.Find()
	if sortFilter == string(model.SortName) {
		opts.SetSort(bson.M{"name": 1})
	} else if sortFilter == string(model.SortRecent) {
		opts.SetSort(bson.M{"updated_at": -1})
	} else {
		opts.SetSort(bson.M{"price": 1})
	}

	cursor, err := collection.Find(ctx, filter, opts)
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

func (s *menuService) CreateMenuItem(item model.MenuItem) (model.MenuItem, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionMenu)
	item.ID = bson.NewObjectID()

	now := time.Now().UnixMilli()
	item.CreatedAt = now
	item.UpdatedAt = now

	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		return model.MenuItem{}, err
	}
	return item, nil
}

func (s *menuService) UpdateMenuItem(id string, updates map[string]interface{}) (model.MenuItem, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.MenuItem{}, errors.New("invalid menu item ID")
	}

	delete(updates, "id")
	delete(updates, "_id")
	delete(updates, "created_at")

	updates["updated_at"] = time.Now().UnixMilli()

	if catStr, ok := updates["categoryId"].(string); ok {
		catID, err := bson.ObjectIDFromHex(catStr)
		if err == nil {
			updates["categoryId"] = catID
		}
	}

	collection := s.db.Collection(config.CollectionMenu)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedItem model.MenuItem
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updates},
		opts,
	).Decode(&updatedItem)

	if err != nil {
		return model.MenuItem{}, err
	}

	return updatedItem, nil
}
