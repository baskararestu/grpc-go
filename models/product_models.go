package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Price    float64            `bson:"price"`
	Category string             `bson:"category"`
}
