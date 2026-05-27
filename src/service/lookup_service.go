package service

import (
	"flamespot-api/src/config"
	"flamespot-api/src/model"
	"flamespot-api/src/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LookupService interface {
	GetActiveCategories() ([]model.CategoryLookup, error)
	GetConfigEnums() (model.EnumConfig, error)
}

type lookupService struct {
	db *mongo.Database
}

func NewLookupService(db *mongo.Database) LookupService {
	return &lookupService{db: db}
}

func (s *lookupService) GetActiveCategories() ([]model.CategoryLookup, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionCategories)

	opts := options.Find().SetSort(bson.M{"name": 1})

	cursor, err := collection.Find(ctx, bson.M{"isActive": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cats []model.CategoryLookup
	if err := cursor.All(ctx, &cats); err != nil {
		return nil, err
	}

	if cats == nil {
		cats = []model.CategoryLookup{}
	}

	return cats, nil
}

func (s *lookupService) GetConfigEnums() (model.EnumConfig, error) {
	return model.EnumConfig{
		Filters: model.FilterConfig{
			Orders: model.OrderFilters{
				Status: model.FilterOrderStatuses(),
				Type:   model.FilterOrderTypes(),
			},
		},
		Sorts: model.SortConfig{
			Categories: model.CategorySorts(),
			Menu:       model.MenuSorts(),
			Orders:     model.OrderSorts(),
		},
	}, nil
}
