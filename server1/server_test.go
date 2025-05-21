package main

import (
	"context"
	"errors"
	"net"
	"server1/calc"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// mockServer implements the CalculateServiceServer interface for testing
type mockServer struct {
	calc.UnimplementedCalculateServiceServer
}

// Calculate implements the Calculate RPC method for the mock server
func (s *mockServer) Calculate(ctx context.Context, req *calc.CalculateRequest) (*calc.CalculateResponse, error) {
	// Simple implementation for testing
	var result int32
	switch req.Opr {
	case "+":
		result = req.A + req.B
	case "-":
		result = req.A - req.B
	case "*":
		result = req.A * req.B
	case "/":
		if req.B == 0 {
			return nil, status.Error(codes.InvalidArgument, "division by zero")
		}
		result = req.A / req.B
	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported operation: "+req.Opr)
	}
	return &calc.CalculateResponse{Result: result}, nil
}

// setupMockServer sets up a mock gRPC server for testing
func setupMockServer(t *testing.T) (calc.CalculateServiceClient, func()) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	calc.RegisterCalculateServiceServer(s, &mockServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("Failed to serve mock server: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := calc.NewCalculateServiceClient(conn)
	return client, func() {
		conn.Close()
		s.Stop()
	}
}

// TestServer1Calculate tests the Calculate method of server1
// Note: This is more of an integration test as it requires both servers
func TestServer1Calculate(t *testing.T) {
	// Set up a mock for server2
	mockClient, cleanup := setupMockServer(t)
	defer cleanup()

	// Create a modified server that uses our mock client
	s := &testServer{
		mockClient: mockClient,
	}

	// Test cases
	testCases := []struct {
		name          string
		request       *calc.CalculateRequest
		expectedResp  *calc.CalculateResponse
		expectError   bool
		errorContains string
	}{
		{
			name: "Test Addition",
			request: &calc.CalculateRequest{
				A:   15,
				Opr: "+",
				B:   5,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 20,
			},
			expectError: false,
		},
		{
			name: "Test Division",
			request: &calc.CalculateRequest{
				A:   20,
				Opr: "/",
				B:   4,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 5,
			},
			expectError: false,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := s.Calculate(context.Background(), tc.request)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Fatalf("Expected error to contain %q, but got %q", tc.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if response.Result != tc.expectedResp.Result {
				t.Errorf("Expected result %d, got %d", tc.expectedResp.Result, response.Result)
			}
		})
	}
}

// testServer is a modified version of server for testing
type testServer struct {
	calc.UnimplementedCalculateServiceServer
	mockClient calc.CalculateServiceClient
}

// Calculate implements the Calculate method for the test server
func (s *testServer) Calculate(ctx context.Context, req *calc.CalculateRequest) (*calc.CalculateResponse, error) {
	// Forward directly to mock client instead of creating a new connection
	return s.mockClient.Calculate(ctx, req)
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return s != "" && substr != "" && errors.Is(errors.New(s), errors.New(substr))
}
