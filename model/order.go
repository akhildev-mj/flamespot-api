package model

type CartItem struct {
	MenuItem MenuItem `json:"menuItem" bson:"menuItem"`
	Quantity int      `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID        string     `json:"id" bson:"_id"`
	TableName string     `json:"tableName" bson:"tableName"`
	Items     []CartItem `json:"items" bson:"items"`
	Total     float64    `json:"total" bson:"total"`
	Status    string     `json:"status" bson:"status"`
	Timestamp int64      `json:"timestamp" bson:"timestamp"`
}
