package main

import (
	"context"
	"log"

	pb "server/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)
	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	log.Printf("Response from server: %s", response.Message)
}
