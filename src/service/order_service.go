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

type OrderService interface {
	GetOrders(limit int64, cursor int64, statusFilter string, typeFilter string, sortFilter string) ([]model.Order, error)
	CreateOrder(order model.Order) (model.Order, error)
	UpdateOrder(id string, updates map[string]interface{}) (model.Order, error)
}

type orderService struct {
	db *mongo.Database
}

func NewOrderService(db *mongo.Database) OrderService {
	return &orderService{db: db}
}

func (s *orderService) GetOrders(limit int64, cursor int64, statusFilter string, typeFilter string, sortFilter string) ([]model.Order, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionOrders)
	filter := bson.M{}

	if cursor > 0 {
		filter["updated_at"] = bson.M{"$lt": cursor}
	}

	if statusFilter != "" && statusFilter != string(model.StatusAll) {
		filter["status"] = model.Status(statusFilter)
	}

	if typeFilter != "" && typeFilter != string(model.TypeAll) {
		filter["type"] = model.Type(typeFilter)
	}

	opts := options.Find().SetLimit(limit)

	if sortFilter == string(model.SortRecent) {
		opts.SetSort(bson.M{"updated_at": -1})
	} else if sortFilter == string(model.SortPrice) {
		opts.SetSort(bson.M{"subTotal": 1})
	} else {
		opts.SetSort(bson.M{"orderedAt": -1})
	}

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

func (s *orderService) CreateOrder(order model.Order) (model.Order, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collection := s.db.Collection(config.CollectionOrders)

	now := time.Now().UnixMilli()

	order.SubTotal = order.Total - order.Discount
	order.CreatedAt = now
	order.UpdatedAt = now

	if order.Status == "" {
		order.Status = model.StatusOrdered
	}

	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order, opts)

	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (s *orderService) UpdateOrder(id string, updates map[string]interface{}) (model.Order, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	if id == "" {
		return model.Order{}, errors.New("invalid order ID")
	}

	delete(updates, "id")
	delete(updates, "_id")
	delete(updates, "created_at")

	collection := s.db.Collection(config.CollectionOrders)

	var existingOrder model.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&existingOrder)
	if err != nil {
		return model.Order{}, errors.New("order not found")
	}

	updates["updated_at"] = time.Now().UnixMilli()

	_, hasTotal := updates["total"]
	_, hasDiscount := updates["discount"]
	_, hasSubTotal := updates["subTotal"]

	if (hasTotal || hasDiscount) && !hasSubTotal {
		newTotal := existingOrder.Total
		if val, ok := updates["total"].(float64); ok {
			newTotal = val
		}

		newDiscount := existingOrder.Discount
		if val, ok := updates["discount"].(float64); ok {
			newDiscount = val
		}

		updates["subTotal"] = newTotal - newDiscount
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedOrder model.Order
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updates},
		opts,
	).Decode(&updatedOrder)

	if err != nil {
		return model.Order{}, err
	}
	return updatedOrder, nil
}
