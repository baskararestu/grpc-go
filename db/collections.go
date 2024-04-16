package db

import "go.mongodb.org/mongo-driver/mongo"

type ColsType struct {
	Products *mongo.Collection
}
