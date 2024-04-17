package services

// import (
// 	"fmt"
// 	"log"

// 	"github.com/baskararestu/grpc-go/db"
// 	m "github.com/baskararestu/grpc-go/models"
// 	"github.com/baskararestu/grpc-go/pb"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	"golang.org/x/net/context"
// )

// type ProductServiceServer struct {
// 	pb.UnimplementedProductServiceServer
// }

// func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
// 	product := req

// 	data := m.ProductItem{
// 		Name:     product.GetName(),
// 		Price:    float64(product.GetPrice()),
// 		Category: product.GetCategory(),
// 	}

// 	result, err := db.GetMongoDB().Products.InsertOne(ctx, data)

// 	// result, err := db.GetMongoDB().Products.InsertOne(context.Background(), data)
// 	if err != nil {
// 		log.Printf("Error inserting product into database: %v", err)
// 		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error inserting product into database: %v", err)}, err
// 	}

// 	oid := result.InsertedID.(primitive.ObjectID)
// 	// if !ok {
// 	// 	log.Printf("Failed to convert InsertedID to ObjectID")
// 	// 	return &pb.ProductResponse{Success: false, Message: "Failed to convert InsertedID to ObjectID"}, nil
// 	// }

// 	// product. = oid.Hex()
// 	return &pb.ProductResponse{Success: true, Message: "Product successfully created", Product: &pb.Product{
// 		Id:       oid.Hex(),
// 		Name:     data.Name,
// 		Price:    float32(data.Price),
// 		Category: data.Category,
// 	}}, nil
// 	// return &pb.ProductResponse{
// 	// 	Success: true,
// 	// 	Message: "Product created successfully",
// 	// 	Product: &pb.Product{
// 	// 		Id:       oid.Hex(),
// 	// 		Name:     product.GetName(),
// 	// 		Price:    float32(product.GetPrice()),
// 	// 		Category: product.GetCategory(),
// 	// 	},
// 	// }, nil
// }

// func (s *ProductServiceServer) ListProduct(req *pb.Empty, stream pb.ProductService_ListProductServer) error {
// 	data := &m.ProductItem{}
// 	ctx := context.Background()
// 	cursor, err := db.GetMongoDB().Products.Find(ctx, bson.M{})
// 	if err != nil {
// 		log.Printf("Error querying products from database: %v", err)
// 		return stream.Send(&pb.ListProductResponse{
// 			Success: false,
// 			Message: fmt.Sprintf("Error querying products from database: %v", err),
// 		})
// 	}
// 	defer cursor.Close(ctx)

// 	for cursor.Next(ctx) {
// 		if err := cursor.Decode(data); err != nil {
// 			log.Printf("Error decoding product: %v", err)
// 			return stream.Send(&pb.ListProductResponse{
// 				Success: false,
// 				Message: fmt.Sprintf("Error decoding product: %v", err),
// 			})
// 		}
// 		if err := stream.Send(&pb.ListProductResponse{
// 			Success: true,
// 			Message: "Products listed successfully",
// 			Products: &pb.Product{
// 				Id:       data.ID.Hex(),
// 				Name:     data.Name,
// 				Price:    float32(data.Price),
// 				Category: data.Category,
// 			},
// 		}); err != nil {
// 			log.Printf("Error sending product response: %v", err)
// 			return err
// 		}
// 	}

// 	if err := cursor.Err(); err != nil {
// 		log.Printf("Unknown cursor error: %v", err)
// 		return stream.Send(&pb.ListProductResponse{
// 			Success: false,
// 			Message: fmt.Sprintf("Error querying products from database: %v", err),
// 		})
// 	}

// 	return nil
// }

// func (s *ProductServiceServer) UpdateProduct(ctx context.Context, req *pb.ProductRequest) (*pb.ProductResponse, error) {
// 	if req.Product.GetId() == "" {
// 		return &pb.ProductResponse{Success: false, Message: "Product ID is required for update"}, nil
// 	}

// 	product := req.GetProduct()
// 	oid, err := primitive.ObjectIDFromHex(product.GetId())
// 	if err != nil {
// 		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Invalid product ID: %s", product.GetId())}, nil
// 	}
// 	update := bson.M{
// 		"name":     product.GetName(),
// 		"price":    float64(product.GetPrice()),
// 		"category": product.GetCategory(),
// 	}

// 	filter := bson.M{
// 		"_id": oid,
// 	}

// 	result := db.DBCollections.Products.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

// 	decoded := m.ProductItem{}
// 	err = result.Decode(&decoded)
// 	if err != nil {
// 		log.Printf("Error decoding product: %v", err)
// 		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error decoding product: %v", err)}, err
// 	}

// 	successMessage := fmt.Sprintf("Product with ID %s updated successfully", product.GetId())
// 	return &pb.ProductResponse{Success: true, Message: successMessage, Product: product}, nil
// }

// func (s *ProductServiceServer) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.DeleteProductResponse, error) {
// 	oid, err := primitive.ObjectIDFromHex(req.GetId())
// 	if err != nil {
// 		return &pb.DeleteProductResponse{Success: false, Message: fmt.Sprintf("Invalid product ID: %s", req.GetId())}, nil
// 	}

// 	_, err = db.GetMongoDB().Products.DeleteOne(ctx, bson.M{"_id": oid})
// 	if err != nil {
// 		return &pb.DeleteProductResponse{Success: false, Message: fmt.Sprintf("Error deleting product: %v", err)}, err
// 	}

// 	return &pb.DeleteProductResponse{Success: true, Message: fmt.Sprintf("Product with ID: %v", req.GetId())}, nil
// }

// func (s *ProductServiceServer) GetProductById(ctx context.Context, req *pb.ProductID) (*pb.ProductResponse, error) {
// 	if req.GetId() == "" {
// 		return &pb.ProductResponse{Success: false, Message: "Product ID is required"}, nil
// 	}

// 	oid, err := primitive.ObjectIDFromHex(req.GetId())
// 	if err != nil {
// 		return &pb.ProductResponse{Success: false, Message: "Product not found"}, nil
// 	}

// 	result := db.GetMongoDB().Products.FindOne(ctx, bson.M{"_id": oid})
// 	data := m.ProductItem{}

// 	if err := result.Decode(&data); err != nil {
// 		return &pb.ProductResponse{Success: false, Message: fmt.Sprintf("Error decoding product: %v", err)}, err
// 	}

// 	response := &pb.ProductResponse{
// 		Success: true,
// 		Message: "Product retrieved successfully",
// 		Product: &pb.Product{
// 			Id:       oid.Hex(),
// 			Name:     data.Name,
// 			Price:    float32(data.Price),
// 			Category: data.Category,
// 		},
// 	}
// 	return response, nil
// }

// func (s *ProductServiceServer) GetRandomProducts(ctx context.Context, req *pb.NRequest) (*pb.ListProductResponse, error) {
// 	cursor, err := db.GetMongoDB().Products.Find(ctx, bson.M{})
// 	if err != nil {
// 		log.Printf("Error querying products from database: %v", err)
// 		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error querying products from database: %v", err)}, err
// 	}
// 	defer cursor.Close(ctx)

// 	var products []*pb.Product
// 	for cursor.Next(ctx) {
// 		var product pb.Product
// 		if err := cursor.Decode(&product); err != nil {
// 			log.Printf("Error decoding product: %v", err)
// 			return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Error decoding product: %v", err)}, err
// 		}
// 		products = append(products, &product)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		log.Printf("Cursor error: %v", err)
// 		return &pb.ListProductResponse{Success: false, Message: fmt.Sprintf("Cursor error: %v", err)}, err
// 	}

// 	rand.NewSource(time.Now().UnixNano())
// 	rand.Shuffle(len(products), func(i, j int) {
// 		products[i], products[j] = products[j], products[i]
// 	})

// 	var productList []*pb.Product
// 	batchSize := int(req.Size)
// 	for i, p := range products {
// 		if i >= batchSize {
// 			break
// 		}
// 		product := &pb.Product{
// 			Id:       p.GetId(),
// 			Name:     p.GetName(),
// 			Price:    p.GetPrice(),
// 			Category: p.GetCategory(),
// 		}
// 		productList = append(productList, product)
// 	}

// 	return &pb.ListProductResponse{
// 		Success:  true,
// 		Message:  fmt.Sprintf("Fetched %d products successfully", len(productList)),
// 		Products: productList,
// 	}, nil
// }
