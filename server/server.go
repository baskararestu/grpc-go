package server

import (
	"context"
	"log"

	"github.com/baskararestu/grpc-go/mapper"
	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"github.com/baskararestu/grpc-go/store"
	"github.com/baskararestu/grpc-go/store/product"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server implements the ProductServiceServer interface.
// It handles the gRPC requests and delegates the actual processing to
// the corresponding functions in the product package.
type server struct {
	pbProduct.UnimplementedProductServiceServer
	GrpcSrv *grpc.Server
	db      *store.MongoDb
}

// type ProductServiceServer struct {
// 	pbProduct.UnimplementedProductServiceServer
// }

// New creates a new instance of the server with the provided database client.
// It sets up the gRPC server, registers the product database service,
// and initializes reflection for gRPC server debugging.
func New(db *store.MongoDb) *server {
	grpcServer := grpc.NewServer()
	srv := &server{
		GrpcSrv: grpcServer,
		db:      db,
	}
	pbProduct.RegisterProductServiceServer(grpcServer, srv)
	reflection.Register(grpcServer)
	return srv
}

// CreateProduct creates a new product in the database.
// It delegates the actual creation logic to the product package's Create function.
func (s *server) CreateProduct(ctx context.Context, in *pbProduct.CreateProductRequest) (*pbProduct.ProductResponse, error) {
	newProduct, err := mapper.CreateProductProtobufToProductModel(in)
	if err != nil {
		return nil, err
	}
	createdProduct, err := product.Create(ctx, s.db, newProduct)
	if err != nil {
		return nil, err
	}
	protoResponse, err := mapper.RespProductModelToProductProtobuf(createdProduct)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// GetProduct retrieves a product by its ID from the database.
// It delegates the actual retrieval logic to the product package's Get function.
func (s *server) GetProductById(ctx context.Context, in *pbProduct.ProductID) (*pbProduct.ProductResponse, error) {
	product, err := product.GetById(ctx, s.db, in)
	if err != nil {
		return nil, errors.Wrapf(err, "getting product with id %s", in.GetId())
	}
	log.Println(product)
	protoResponse, err := mapper.ResProductModelToProductProtobuf(product)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// UpdateProduct updates an existing product in the database.
// It delegates the actual update logic to the product package's Update function.
func (s *server) UpdateProduct(ctx context.Context, in *pbProduct.ProductRequest) (*pbProduct.ProductResponse, error) {
	productToUpdate, err := mapper.ReqProductProtobufToProductModel(in)
	if err != nil {
		return nil, err
	}
	updatedProduct, err := product.Update(ctx, s.db, productToUpdate)
	if err != nil {
		return nil, err
	}
	protoResponse, err := mapper.RespProductModelToProductProtobuf(updatedProduct)
	if err != nil {
		return nil, err
	}
	return protoResponse, nil
}

// DeleteProduct deletes a product from the database.
// It delegates the actual deletion logic to the product package's Delete function.
func (s *server) DeleteProduct(ctx context.Context, in *pbProduct.ProductID) (*pbProduct.DeleteProductResponse, error) {
	resp, err := product.Delete(ctx, s.db, in)
	if err != nil {
		return nil, errors.Wrapf(err, "deleting product with uuid %s", in.GetId())
	}
	return resp, nil
}

// ListProduct streams all the products in the database to the client.
func (s *server) ListProduct(in *pbProduct.Empty, stream pbProduct.ProductService_ListProductServer) error {
	products, err := product.List(context.Background(), s.db, in)
	if err != nil {
		return errors.Wrap(err, "listing products")
	}

	response := &pbProduct.ListProductResponse{
		Success:  true,
		Message:  "Products retrieved successfully",
		Products: []*pbProduct.Product{},
	}

	for _, dbProduct := range products {
		protoProduct, err := mapper.ProductModelToProductProtobuf(dbProduct)
		if err != nil {
			return err
		}
		response.Products = append(response.Products, protoProduct)
	}
	if err := stream.Send(response); err != nil {
		return errors.Wrap(err, "sending product list over stream")
	}
	return nil
}

// GetRandomProducts streams a random number of products from the database to the client.
func (s *server) GetRandomProducts(req *pbProduct.NRequest, stream pbProduct.ProductService_GetRandomProductsServer) error {
	products, err := product.ListRandom(context.Background(), s.db, req)
	if err != nil {
		return errors.Wrap(err, "listing random products")
	}

	response := &pbProduct.ListProductResponse{
		Success:  true,
		Message:  "Random Products retrieved successfully",
		Products: []*pbProduct.Product{},
	}

	for _, dbProduct := range products {
		protoProduct, err := mapper.ProductModelToProductProtobuf(dbProduct)
		if err != nil {
			return err
		}
		response.Products = append(response.Products, protoProduct)
	}
	if err := stream.Send(response); err != nil {
		return errors.Wrap(err, "sending random product list over stream")
	}
	return nil
}
