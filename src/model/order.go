package model

type OrderMenuItem struct {
	ID         string  `json:"id" bson:"id"`
	CategoryID string  `json:"categoryId" bson:"categoryId"`
	Name       string  `json:"name" bson:"name"`
	Price      float64 `json:"price" bson:"price"`
	URL        string  `json:"url" bson:"url"`
	IsActive   bool    `json:"isActive" bson:"isActive"`
}

type CartItem struct {
	MenuItem OrderMenuItem `json:"menuItem" bson:"menuItem"`
	Quantity int           `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID        string     `json:"id" bson:"_id"`
	Type      Type       `json:"type" bson:"type"`
	TableId   int        `json:"tableId" bson:"tableId"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"total"`
	Discount  float64    `json:"discount" bson:"discount"`
	SubTotal  float64    `json:"subTotal" bson:"subTotal"`
	Status    Status     `json:"status" bson:"status"`
	OrderedAt int64      `json:"orderedAt" bson:"orderedAt"`
	BilledAt  int64      `json:"billedAt" bson:"billedAt"`
	SavedAt   int64      `json:"savedAt" bson:"savedAt"`
	DeletedAt int64      `json:"deletedAt" bson:"deletedAt"`
	CreatedAt int64      `json:"-" bson:"created_at"`
	UpdatedAt int64      `json:"-" bson:"updated_at"`
}
