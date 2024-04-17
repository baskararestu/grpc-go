package main

import (
	"context"
	"log"
	"strconv"
	"time"

	pbProduct "github.com/baskararestu/grpc-go/pb/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pbProduct.NewProductServiceClient(conn)

	receivedTimes := make([]time.Time, 0)

	for i := 0; i < 20; i++ {
		product := &pbProduct.CreateProductRequest{
			Name:     "Product " + strconv.Itoa(i+1),
			Price:    float32((i + 1) * 10),
			Category: "Category " + strconv.Itoa(i+1),
		}

		resp, err := client.CreateProduct(context.Background(), product)
		if err != nil {
			log.Fatalf("Failed to create product: %v", err)
		}

		receivedTimes = append(receivedTimes, time.Now())

		log.Printf("Received product response: %v", resp)
	}

	isSorted := true
	for i := 1; i < len(receivedTimes); i++ {
		if receivedTimes[i].Before(receivedTimes[i-1]) {
			isSorted = false
			break
		}
	}

	if isSorted {
		log.Println("Data received in random order.")
	} else {
		log.Println("Data received in sequential order.")
	}

	timeDifferences := make([]float64, 0)
	for i := 1; i < len(receivedTimes); i++ {
		diff := receivedTimes[i].Sub(receivedTimes[i-1]).Seconds()
		timeDifferences = append(timeDifferences, diff)
	}

	log.Printf("Time differences between received data: %v", timeDifferences)
}
