package services

import (
	"fmt"
	"log"

	"github.com/baskararestu/grpc-go/db"
	"github.com/baskararestu/grpc-go/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"

	m "github.com/baskararestu/grpc-go/models"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := &m.ProductItem{
		Name:     req.GetName(),
		Price:    float64(req.GetPrice()),
		Category: req.GetCategory(),
	}

	count, err := db.GetMongoDB().Products.CountDocuments(ctx, bson.M{"name": req.Name})
	if err != nil {
		log.Printf("Error checking for existing product with name %s: %v", req.Name, err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error checking for existing product with name %s: %v", req.Name, err)}, err
	}

	if count > 0 {
		log.Printf("Product with name %s already exists", req.Name)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Product with name %s already exists", req.Name)}, nil
	}

	result, err := db.GetMongoDB().Products.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error inserting product into database: %v", err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error inserting product into database: %v", err)}, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting inserted ID to ObjectID")
		return &pb.ProductResponse{Success: false, Message: "Error converting inserted ID to ObjectID"}, err
	}

	product.ID = insertedID

	resp := &pb.ProductResponse{
		Success: true,
		Message: "Product created successfully",
		Product: []*pb.Product{{
			Id:       product.ID.Hex(),
			Name:     product.Name,
			Price:    float32(product.Price),
			Category: product.Category,
		}},
	}

	return resp, nil
}
