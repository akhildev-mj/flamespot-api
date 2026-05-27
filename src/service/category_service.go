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

type CategoryService interface {
	GetCategories(isActiveFilter *bool, sortFilter string) ([]model.Category, error)
	CreateCategory(cat model.Category) (model.Category, error)
	UpdateCategory(id string, updates map[string]interface{}) (model.Category, error)
}

type categoryService struct {
	db *mongo.Database
}

func NewCategoryService(db *mongo.Database) CategoryService {
	return &categoryService{db: db}
}

func (s *categoryService) GetCategories(isActiveFilter *bool, sortFilter string) ([]model.Category, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionCategories)

	filter := bson.M{}
	if isActiveFilter != nil {
		filter["isActive"] = *isActiveFilter
	}

	opts := options.Find()
	if sortFilter == string(model.SortRecent) {
		opts.SetSort(bson.M{"updated_at": -1})
	} else {
		opts.SetSort(bson.M{"name": 1})
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cats []model.Category
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, err
	}

	if cats == nil {
		cats = []model.Category{}
	}
	return cats, nil
}

func (s *categoryService) CreateCategory(cat model.Category) (model.Category, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionCategories)
	cat.ID = bson.NewObjectID()

	now := time.Now().UnixMilli()
	cat.CreatedAt = now
	cat.UpdatedAt = now

	_, err := collection.InsertOne(ctx, cat)
	if err != nil {
		return model.Category{}, err
	}
	return cat, nil
}

func (s *categoryService) UpdateCategory(id string, updates map[string]interface{}) (model.Category, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.Category{}, errors.New("invalid ID")
	}

	delete(updates, "id")
	delete(updates, "_id")
	delete(updates, "created_at")

	updates["updated_at"] = time.Now().UnixMilli()

	collection := s.db.Collection(config.CollectionCategories)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated model.Category
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updates},
		opts,
	).Decode(&updated)

	if err != nil {
		return model.Category{}, err
	}
	return updated, nil
}
