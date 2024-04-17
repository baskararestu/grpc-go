// Package mapper provides functions for converting between Protobuf messages
// and MongoDB models in the context of a product database.
// The functions in this package handle the conversion of product data between
// the Protobuf representation used in the API and the MongoDB model representation
// used in the data store.
package mapper

import (
	"github.com/baskararestu/grpc-go/models"
	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductProtobufToProductModel converts a Protobuf Product message to a MongoDB Product model.
// func ProductProtobufToProductModel(product *pbProduct.Product) (*models.ProductItem, error) {
// 	// Convert string ID to primitive.ObjectID
// 	objectID, err := primitive.ObjectIDFromHex(product.GetId())
// 	if err != nil {
// 		return nil, err
// 	}

// 	dbProduct := &models.ProductItem{
// 		ID:       objectID,
// 		Name:     product.Name,
// 		Price:    float64(product.GetPrice()),
// 		Category: product.GetCategory(),
// 	}
// 	return dbProduct, nil
// }

func CreateProductProtobufToProductModel(product *pbProduct.CreateProductRequest) (*models.ProductItem, error) {
	dbProduct := &models.ProductItem{
		Name:     product.Name,
		Price:    float64(product.GetPrice()),
		Category: product.GetCategory(),
	}
	return dbProduct, nil
}

// ProductModelToProductProtobuf converts a MongoDB Product model to a Protobuf Product message.
func ProductModelToProductProtobuf(dbProduct *models.ProductItem) (*pbProduct.Product, error) {
	product := &pbProduct.Product{
		Id:       dbProduct.ID.Hex(),
		Name:     dbProduct.Name,
		Price:    float32(dbProduct.Price),
		Category: dbProduct.Category,
	}
	return product, nil
}

func RespProductModelToProductProtobuf(dbProduct *models.ProductItem) (*pbProduct.ProductResponse, error) {
	product := &pbProduct.Product{
		Id:       dbProduct.ID.Hex(),
		Name:     dbProduct.Name,
		Price:    float32(dbProduct.Price),
		Category: dbProduct.Category,
	}
	response := &pbProduct.ProductResponse{
		Success: true,
		Message: "Product created successfully",
		Product: []*pbProduct.Product{product},
	}
	return response, nil
}

func ResProductModelToProductProtobuf(dbProduct *models.ProductItem) (*pbProduct.ProductResponse, error) {
	product := &pbProduct.Product{
		Id:       dbProduct.ID.Hex(),
		Name:     dbProduct.Name,
		Price:    float32(dbProduct.Price),
		Category: dbProduct.Category,
	}

	response := &pbProduct.ProductResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Product: []*pbProduct.Product{product},
	}
	return response, nil
}

func ReqProductProtobufToProductModel(product *pbProduct.ProductRequest) (*models.ProductItem, error) {
	// Convert string ID to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(product.GetProduct().GetId())
	if err != nil {
		return nil, err
	}

	dbProduct := &models.ProductItem{
		ID:       objectID,
		Name:     product.GetProduct().Name,
		Price:    float64(product.GetProduct().GetPrice()),
		Category: product.GetProduct().GetCategory(),
	}
	return dbProduct, nil
}

// ProductModelListToListProductsResponse converts a list of MongoDB Product models to a Protobuf ListProductsResponse message.
// func ProductModelListToListProductsResponse(dbProducts []*models.ProductItem) (*pbProduct.ListProductResponse, error) {
// 	response := &pbProduct.ListProductResponse{}
// 	products := []*pbProduct.Product{}
// 	for _, dbProduct := range dbProducts {
// 		product, err := ProductModelToProductProtobuf(dbProduct)
// 		if err != nil {
// 			return nil, err
// 		}
// 		products = append(products, product)
// 	}
// 	response.Products = products
// 	return response, nil
// }
