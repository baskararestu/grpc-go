package main

import (
	"fmt"

	"log"
	"net"

	"github.com/baskararestu/grpc-go/db"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/baskararestu/grpc-go/pb"
	"github.com/baskararestu/grpc-go/services"
)

func init() {
	fmt.Println("init runs")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("err loading env %v", err)
	}
}

func main() {
	db.ConnectMongo()

	grpcServer := grpc.NewServer()

	p := &services.ProductServiceServer{}

	pb.RegisterProductServiceServer(grpcServer, p)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server started on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
