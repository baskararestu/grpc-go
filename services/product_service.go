package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

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

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
	product := &pb.Product{
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
	}

	count, err := db.GetMongoDB().Products.CountDocuments(ctx, bson.M{"name": product.GetName()})
	if err != nil {
		log.Printf("Error checking for existing product with name %s: %v", product.GetName(), err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error checking for existing product with name %s: %v", req.Name, err)}, err
	}

	if count > 0 {
		log.Printf("Product with name %s already exists", product.GetName())
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Product with name %s already exists", product.GetName())}, nil
	}

	res, err := db.GetMongoDB().Products.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error inserting product into database: %v", err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error inserting product into database: %v", err)}, err
	}

	insertedID := res.InsertedID.(primitive.ObjectID)
	product.Id = insertedID.Hex()

	resp := &pb.ProductResponse{
		Success: true,
		Message: "Product created successfully",
		Product: []*pb.Product{{
			Id:       product.GetId(),
			Name:     product.GetName(),
			Price:    product.GetPrice(),
			Category: product.GetCategory(),
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

	var products []*pb.Product
	for cursor.Next(ctx) {
		var product pb.Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("Error decoding product: %v", err)
			return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error decoding product: %v", err)}, err
		}
		product.Id = product.GetId()
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Cursor error: %v", err)}, err
	}

	return &pb.ListProductResponse{
		Success:  true,
		Message:  "Products listed successfully",
		Products: products,
	}, nil
}

func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.ProductUpdateRequest) (*pb.ProductResponse, error) {
	if req.GetId() == "" {
		return &pb.ProductResponse{Success: false, Message: "Product ID is required for update"}, nil
	}

	filter := bson.M{"_id": req.GetId()}

	update := bson.M{
		"$set": bson.M{
			"name":     req.Product.GetName(),
			"price":    req.Product.GetPrice(),
			"category": req.Product.GetCategory(),
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
	if req.GetId() == "" {
		return &pb.DeleteProductResponse{Success: false, Message: "Product ID is required for deletion"}, nil
	}

	filter := bson.M{"_id": req.GetId()}

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
	if req.GetId() == "" {
		return &pb.ProductResponse{Success: false, Message: "Product ID is required for retrieval"}, nil
	}

	filter := bson.M{"_id": req.GetId()}

	var product *pb.Product
	err := db.GetMongoDB().Products.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Product not found with ID: %s", req.Id)
			return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Product not found with ID: %s", req.Id)}, nil
		}
		log.Printf("Error retrieving product: %v", err)
		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error retrieving product: %v", err)}, err
	}

	return &pb.ProductResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Product: []*pb.Product{{
			Id:       product.GetId(),
			Name:     product.GetName(),
			Price:    product.GetPrice(),
			Category: product.GetCategory(),
		}},
	}, nil
}

func (s *ProductServiceServer) GetRandomProducts(ctx context.Context, req *pb.NRequest) (*pb.ListProductResponse, error) {
	cursor, err := db.GetMongoDB().Products.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error querying products from database: %v", err)
		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error querying products from database: %v", err)}, err
	}
	defer cursor.Close(ctx)

	var products []*pb.Product
	for cursor.Next(ctx) {
		var product pb.Product
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

	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(products), func(i, j int) {
		products[i], products[j] = products[j], products[i]
	})

	var productList []*pb.Product
	batchSize := int(req.Size)
	for i, p := range products {
		if i >= batchSize {
			break
		}
		product := &pb.Product{
			Id:       p.GetId(),
			Name:     p.GetName(),
			Price:    p.GetPrice(),
			Category: p.GetCategory(),
		}
		productList = append(productList, product)
	}

	return &pb.ListProductResponse{
		Success:  true,
		Message:  fmt.Sprintf("Fetched %d products successfully", len(productList)),
		Products: productList,
	}, nil
}
