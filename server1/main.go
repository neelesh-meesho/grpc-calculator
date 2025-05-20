package main

import (
	"calc/calc"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calc.UnimplementedCalculateServiceServer
}

func (s *server) Calculate(ctx context.Context, req *calc.CalculateRequest) (*calc.CalculateResponse, error) {
	a := req.A
	b := req.B
	opr := req.Opr

	var result int32

	switch opr {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	}

	return &calc.CalculateResponse{Result: result}, nil
}

func createServer() {
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	calc.RegisterCalculateServiceServer(s, &server{})

	lis, err := net.Listen("tcp", ":50051")

	fmt.Println("Starting server on port: 50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.Serve(lis)
}

func main() {
	fmt.Println("Welcome to my the calculator!")
	createServer()
}
