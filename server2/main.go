package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"server2/calc"

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
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unsupported operation: %s", opr)
	}

	log.Printf("Calculation performed: %d %s %d = %d", a, opr, b, result)
	return &calc.CalculateResponse{Result: result}, nil
}

func main() {
	port := ":50052"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calc.RegisterCalculateServiceServer(s, &server{})

	fmt.Printf("Server2 (calculation service) started on port%s\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
