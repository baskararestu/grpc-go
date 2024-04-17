package mapper

import (
	"testing"

	"github.com/baskararestu/grpc-go/models"
	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func TestProductProtobufToProductModel(t *testing.T) {
// 	objID := primitive.NewObjectID()
// 	p := &pbProduct.Product{
// 		Id:       objID.Hex(),
// 		Name:     "Product1",
// 		Category: "category",
// 		Price:    10.0,
// 	}

// 	dbProduct, err := ProductProtobufToProductModel(p)
// 	assert.Nil(t, err)
// 	assert.Equal(t, objID.Hex(), dbProduct.ID.Hex())
// 	assert.Equal(t, "Product1", dbProduct.Name)
// 	assert.Equal(t, "category", dbProduct.Category)
// 	assert.Equal(t, 10.0, dbProduct.Price)
// }

func TestProdutcModelToProductProtobuf(t *testing.T) {
	objID := primitive.NewObjectID()
	dbProduct := &models.ProductItem{
		ID:       objID,
		Name:     "name",
		Category: "category",
		Price:    1,
	}

	expectedProduct := &pbProduct.Product{
		Id:       objID.Hex(),
		Name:     "name",
		Category: "category",
		Price:    1,
	}

	product, err := ProductModelToProductProtobuf(dbProduct)
	assert.Nil(t, err)
	assert.Equal(t, expectedProduct, product)
}

// func TestProductModelListToListProductsResponse(t *testing.T) {
// 	objID := primitive.NewObjectID()
// 	dbProducts := []*models.ProductItem{
// 		{
// 			ID:       objID,
// 			Name:     "name",
// 			Category: "category",
// 			Price:    1,
// 		},
// 	}

// 	expectedResponse := &pbProduct.ListProductResponse{
// 		Products: []*pbProduct.Product{
// 			{
// 				Id:       objID.Hex(),
// 				Name:     "name",
// 				Category: "category",
// 				Price:    1,
// 			},
// 		},
// 	}

// 	response, err := ProductModelListToListProductsResponse(dbProducts)
// 	assert.Nil(t, err)
// 	assert.Equal(t, expectedResponse, response)
// }
