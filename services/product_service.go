package services

import (
	"fmt"
	"log"

	"github.com/baskararestu/grpc-go/db"
	"github.com/baskararestu/grpc-go/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
}

type Product struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Price    float32            `bson:"price"`
	Category string             `bson:"category"`
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	product := &Product{
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
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
			Price:    product.Price,
			Category: product.Category,
		}},
	}

	return resp, nil
}

func (s *ProductServiceServer) ListProduct(ctx context.Context, req *pb.Empty) (*pb.ListProductResponse, error) {
	cursor, err := db.GetMongoDB().Products.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error querying products from database: %v", err)
		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error querying products from database: %v", err)}, err
	}
	defer cursor.Close(ctx)

	var products []*Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("Error decoding product: %v", err)
			return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error decoding product: %v", err)}, err
		}
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Cursor error: %v", err)}, err
	}

	var productList []*pb.Product
	for _, p := range products {
		product := &pb.Product{
			Id:       p.ID.Hex(),
			Name:     p.Name,
			Price:    p.Price,
			Category: p.Category,
		}
		productList = append(productList, product)
	}

	listProductResponse := &pb.ListProductResponse{
		Success:  true,
		Message:  "Products listed successfully",
		Products: productList,
	}

	return listProductResponse, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.ProductUpdateRequest) (*pb.ProductResponse, error) {
	if req.Id == "" {
		return &pb.ProductResponse{Success: false, Message: "Product ID is required for update"}, nil
	}

	productID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Printf("Invalid product ID format: %v", err)
		return &pb.ProductResponse{Success: false, Message: "Invalid product ID format"}, err
	}

	filter := bson.M{"_id": productID}

	update := bson.M{
		"$set": bson.M{
			"name":     req.Product.Name,
			"price":    req.Product.Price,
			"category": req.Product.Category,
		},
	}

	updateResult, err := db.GetMongoDB().Products.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error updating product: %v", err)}, err
	}

	if updateResult.ModifiedCount == 0 {
		return &pb.ProductResponse{Success: false, Message: "Product not found or no changes were made"}, nil
	}

	successMessage := fmt.Sprintf("Product with ID %s updated successfully", req.Id)
	return &pb.ProductResponse{Success: true, Message: successMessage}, nil
}

func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.DeleteProductResponse, error) {
	if req.Id == "" {
		return &pb.DeleteProductResponse{Success: false, Message: "Product ID is required for deletion"}, nil
	}

	productID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Printf("Invalid product ID format: %v", err)
		return &pb.DeleteProductResponse{Success: false, Message: "Invalid product ID format"}, err
	}

	filter := bson.M{"_id": productID}

	deleteResult, err := db.GetMongoDB().Products.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return &pb.DeleteProductResponse{Success: false, Message: fmt.Sprintf("Error deleting product: %v", err)}, err
	}

	if deleteResult.DeletedCount == 0 {
		return &pb.DeleteProductResponse{Success: false, Message: "Product not found"}, nil
	}

	return &pb.DeleteProductResponse{Success: true, Message: "Product deleted successfully"}, nil
}

func (s *ProductServiceServer) GetProductById(ctx context.Context, req *pb.ProductID) (*pb.ProductResponse, error) {
	if req.Id == "" {
		return &pb.ProductResponse{Success: false, Message: "Product ID is required for retrieval"}, nil
	}

	productID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		log.Printf("Invalid product ID format: %v", err)
		return &pb.ProductResponse{Success: false, Message: "Invalid product ID format"}, err
	}

	filter := bson.M{"_id": productID}

	var product Product
	err = db.GetMongoDB().Products.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Product not found with ID: %s", req.Id)
			return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Product not found with ID: %s", req.Id)}, nil
		}
		log.Printf("Error retrieving product: %v", err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error retrieving product: %v", err)}, err
	}

	productResponse := &pb.ProductResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Product: []*pb.Product{{
			Id:       product.ID.Hex(),
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
		}},
	}

	return productResponse, nil
}
