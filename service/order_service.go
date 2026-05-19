package service

import (
	"flamespot-api/config"
	"flamespot-api/model"
	"flamespot-api/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OrderService interface {
	CreateOrder(order model.Order) error
	GetOrders(limit int64, cursor int64) ([]model.Order, error)
}

type orderService struct {
	db *mongo.Database
}

func NewOrderService(db *mongo.Database) OrderService {
	return &orderService{db: db}
}

func (s *orderService) CreateOrder(order model.Order) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionOrders)
	opts := options.Replace().SetUpsert(true)

	_, err := collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order, opts)
	return err
}

func (s *orderService) GetOrders(limit int64, cursor int64) ([]model.Order, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionOrders)

	filter := bson.M{}
	if cursor > 0 {
		filter = bson.M{"timestamp": bson.M{"$lt": cursor}}
	}

	opts := options.Find().SetSort(bson.M{"timestamp": -1}).SetLimit(limit)

	cursorResult, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursorResult.Close(ctx)

	var orders []model.Order
	if err := cursorResult.All(ctx, &orders); err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []model.Order{}
	}

	return orders, nil
}
