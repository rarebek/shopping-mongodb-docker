package models

type User struct {
	// ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID string  `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Name string  `json:"name,omitempty" bson:"name,omitempty"`
	Age  float64 `json:"age,omitempty" bson:"age,omitempty"`
}

type Product struct {
	// ObjectID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UUID     string  `json:"uuid,omitempty" bson:"uuid,omitempty"`
	Name     string  `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64 `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int     `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type Cart struct {
	// ObjectID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserUUID    string    `json:"user_uuid,omitempty" bson:"user_uuid,omitempty"`
	Products    []Product `json:"product_uuids,omitempty" bson:"product_uuids,omitempty"`
	TotalAmount float64   `json:"total_amount,omitempty" bson:"total_amount,omitempty"`
}
