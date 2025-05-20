package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"server1/calc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	calc.UnimplementedCalculateServiceServer
}

func (s *server) Calculate(ctx context.Context, req *calc.CalculateRequest) (*calc.CalculateResponse, error) {
	log.Printf("Server1 received calculation request: %d %s %d", req.A, req.Opr, req.B)

	// Connect to server2
	server2Addr := "localhost:50052"
	conn, err := grpc.DialContext(ctx, server2Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server2: %v", err)
		return nil, fmt.Errorf("failed to connect to calculation service: %v", err)
	}
	defer conn.Close()

	// Create client
	client := calc.NewCalculateServiceClient(conn)

	// Set timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Forward the request to server2
	response, err := client.Calculate(ctx, req)
	if err != nil {
		log.Printf("Error from server2: %v", err)
		return nil, err
	}

	log.Printf("Server1 received result from server2: %d", response.Result)
	return response, nil
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calc.RegisterCalculateServiceServer(s, &server{})

	fmt.Printf("Server1 started on port%s, will delegate calculations to server2\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
