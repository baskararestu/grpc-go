package product

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/baskararestu/grpc-go/models"
	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"github.com/baskararestu/grpc-go/store"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
)

func TestCreate(t *testing.T) {
	objectID := primitive.NewObjectID()
	testCases := []struct {
		name                     string
		input                    *models.ProductItem
		mockInsertIntoCollection func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error)
		expectedOutput           *models.ProductItem
		expectedError            error
	}{
		{
			name: "happy path",
			input: &models.ProductItem{
				Name:     "name",
				Price:    1,
				Category: "category",
			},
			expectedOutput: &models.ProductItem{
				ID:       objectID,
				Name:     "name",
				Price:    1,
				Category: "category",
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return &mongo.InsertOneResult{InsertedID: objectID}, nil
			},
		},
		{
			name: "error",
			input: &models.ProductItem{
				Name:     "name",
				Price:    1,
				Category: "category",
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("inserting product: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			insertIntoCollection = tc.mockInsertIntoCollection
			output, err := Create(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				tc.expectedOutput.ID = objectID
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestGet(t *testing.T) {
	objectID := primitive.NewObjectID()

	testCases := []struct {
		name           string
		mockFindOne    func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.ProductItem) error
		expectedOutput *models.ProductItem
		expectedError  error
	}{
		{
			name: "happy path",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.ProductItem) error {

				p.ID = objectID
				p.Name = "name"
				p.Category = "Category"
				p.Price = 1
				return nil
			},
			expectedOutput: &models.ProductItem{
				ID:       objectID,
				Name:     "name",
				Category: "Category",
				Price:    1,
			},
		},
		{
			name: "error",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.ProductItem) error {
				return errors.New("random error")
			},
			expectedError: errors.New(`getting product with id "` + objectID.Hex() + `": random error`),
		},
		{
			name: "document not found",
			mockFindOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, p *models.ProductItem) error {
				return mongo.ErrNoDocuments
			},
			expectedError: errors.New(`product with id "` + objectID.Hex() + `" does not exist`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			findOne = tc.mockFindOne
			output, err := GetById(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &pbProduct.ProductID{Id: objectID.Hex()})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	objectID := primitive.NewObjectID()

	testCases := []struct {
		name           string
		input          *models.ProductItem
		mockUpdateOne  func(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
		expectedOutput *models.ProductItem
		expectedError  error
	}{
		{
			name: "happy path",
			input: &models.ProductItem{
				Name:     "name",
				Price:    1,
				Category: "category",
			},
			mockUpdateOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
				return &mongo.UpdateResult{}, nil
			},
			expectedOutput: &models.ProductItem{
				ID:       objectID,
				Name:     "name",
				Price:    1,
				Category: "category",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateOne = tc.mockUpdateOne
			output, err := Update(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &models.ProductItem{
				ID:       tc.input.ID,
				Name:     tc.input.Name,
				Price:    tc.input.Price,
				Category: tc.input.Category,
			})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				tc.expectedOutput.ID = output.ID
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	objectID := primitive.NewObjectID()
	testCases := []struct {
		name           string
		mockDeleteOne  func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error)
		expectedOutput *pbProduct.DeleteProductResponse
		expectedError  error
	}{
		{
			name: "happy path",
			mockDeleteOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
				return nil, nil
			},
			expectedOutput: &pbProduct.DeleteProductResponse{
				Success: true,
				Message: fmt.Sprintf(`product with id "%s" deleted`, objectID.Hex()),
			},
		},
		{
			name: "error",
			mockDeleteOne: func(ctx context.Context, collection *mongo.Collection, filter interface{}) (*mongo.DeleteResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New(`deleting product with id "` + objectID.Hex() + `": random error`),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteOne = tc.mockDeleteOne
			output, err := Delete(context.TODO(), &store.MongoDb{DatabaseName: "db", Client: &mongo.Client{}}, &pbProduct.ProductID{Id: objectID.Hex()})
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error %v, got nil", tc.expectedError)
				}
				log.Println(tc.expectedOutput)
				log.Println(output)
				require.True(t, proto.Equal(tc.expectedOutput, output))
			}
		})
	}
}
