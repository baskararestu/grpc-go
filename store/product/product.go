package product

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/baskararestu/grpc-go/models"
	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"github.com/baskararestu/grpc-go/store"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "products"

// Cursor is an interface that defines the methods necessary for iterating
// over query results in a data layer.
// This interface is particularly useful for simplifying unit tests
// by allowing the implementation of mock cursors that can be used
// for testing data retrieval and manipulation operations.
type Cursor interface {
	Decode(interface{}) error
	Err() error
	Close(context.Context) error
	Next(context.Context) bool
}

type cursorWrapper struct {
	*mongo.Cursor
}

// For ease of unit testing.
var (
	insertIntoCollection = func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(ctx, document)
	}
	findOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.ProductItem) error {
		sr := collection.FindOne(ctx, filter)
		return sr.Decode(p)
	}
	find = func(ctx context.Context, collection *mongo.Collection, filter interface{}) (Cursor, error) {
		cur, err := collection.Find(ctx, filter)
		return &cursorWrapper{cur}, err
	}
	updateOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
		return collection.UpdateOne(ctx, filter, update)
	}
	deleteOne = func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
		return collection.DeleteOne(ctx, filter)
	}
)

// Get retrieves a product from the database by id.
func GetById(ctx context.Context, db *store.MongoDb, req *pbProduct.ProductID) (*models.ProductItem, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	var product models.ProductItem
	objectID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, errors.Wrapf(err, `converting "%s" to ObjectID`, req.GetId())
	}
	err = findOne(ctx, coll, bson.M{"_id": objectID}, &product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(`product with id "%s" does not exist`, req.GetId())
		}
		return nil, errors.Wrapf(err, `getting product with id "%s"`, req.GetId())
	}
	return &product, nil
}

// Create creates a new product in the database.
func Create(ctx context.Context, db *store.MongoDb, newProduct *models.ProductItem) (*models.ProductItem, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	result, err := insertIntoCollection(ctx, coll, newProduct)
	if err != nil {
		return nil, errors.Wrap(err, "inserting product")
	}
	newProduct.ID = result.InsertedID.(primitive.ObjectID)
	return newProduct, nil
}

// Update updates a product in the database.
func Update(ctx context.Context, db *store.MongoDb, productToUpdate *models.ProductItem) (*models.ProductItem, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	_, err := updateOne(ctx, coll, bson.M{"_id": productToUpdate.ID}, bson.M{"$set": productToUpdate})
	if err != nil {
		return nil, errors.Wrapf(err, `updating product with id "%s"`, productToUpdate.ID)
	}
	return productToUpdate, nil
}

// Delete deletes a product from the database by id.
func Delete(ctx context.Context, db *store.MongoDb, req *pbProduct.ProductID) (*pbProduct.DeleteProductResponse, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	objectID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, errors.Wrapf(err, `converting "%s" to ObjectID`, req.GetId())
	}

	_, err = deleteOne(ctx, coll, bson.M{"_id": objectID})
	if err != nil {
		return nil, errors.Wrapf(err, `deleting product with id "%s"`, req.GetId())
	}
	return &pbProduct.DeleteProductResponse{
		Success: true,
		Message: fmt.Sprintf(`product with id "%s" deleted`, req.GetId()),
	}, nil
}

// List lists all products in the database.
func List(ctx context.Context, db *store.MongoDb, req *pbProduct.Empty) ([]*models.ProductItem, error) {
	coll := db.Client.Database(db.DatabaseName).Collection(collectionName)
	cur, err := find(ctx, coll, bson.M{})
	if err != nil {
		return nil, errors.Wrap(err, "finding products")
	}
	defer cur.Close(ctx)
	var products []*models.ProductItem
	for cur.Next(ctx) {
		var product models.ProductItem
		if err = cur.Decode(&product); err != nil {
			return nil, errors.Wrap(err, "decoding product")
		}
		products = append(products, &product)
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor error")
	}
	return products, nil
}

// ListRandom retrieves a random number of products from the database.
func ListRandom(ctx context.Context, db *store.MongoDb, req *pbProduct.NRequest) ([]*models.ProductItem, error) {
	// coll := db.Client.Database(db.DatabaseName).Collection(collectionName)

	// Retrieve all products from the database
	allProducts, err := List(ctx, db, &pbProduct.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve products")
	}

	// Get the size of the requested products
	size := int(req.GetSize())

	// If the requested size is greater than the total number of products, return all products
	if size >= len(allProducts) {
		return allProducts, nil
	}

	// Shuffle the products randomly
	rand.Shuffle(len(allProducts), func(i, j int) {
		allProducts[i], allProducts[j] = allProducts[j], allProducts[i]
	})

	// Return a slice of randomly selected products
	return allProducts[:size], nil
}
