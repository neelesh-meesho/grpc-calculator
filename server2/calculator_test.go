package main

import (
	"context"
	"server2/calc"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	// Create server instance
	s := &server{}
	ctx := context.Background()

	// Define test cases
	testCases := []struct {
		name          string
		request       *calc.CalculateRequest
		expectedResp  *calc.CalculateResponse
		expectError   bool
		errorContains string
	}{
		{
			name: "Addition",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "+",
				B:   5,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 15,
			},
			expectError: false,
		},
		{
			name: "Subtraction",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "-",
				B:   5,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 5,
			},
			expectError: false,
		},
		{
			name: "Multiplication",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "*",
				B:   5,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 50,
			},
			expectError: false,
		},
		{
			name: "Division",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "/",
				B:   5,
			},
			expectedResp: &calc.CalculateResponse{
				Result: 2,
			},
			expectError: false,
		},
		{
			name: "Division by zero",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "/",
				B:   0,
			},
			expectedResp:  nil,
			expectError:   true,
			errorContains: "division by zero",
		},
		{
			name: "Unsupported operation",
			request: &calc.CalculateRequest{
				A:   10,
				Opr: "%",
				B:   5,
			},
			expectedResp:  nil,
			expectError:   true,
			errorContains: "unsupported operation",
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := s.Calculate(ctx, tc.request)

			// Check error condition
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Fatalf("Expected error to contain %q, but got %q", tc.errorContains, err.Error())
				}
				return
			}

			// Check no error happened
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result matches expected
			if response.Result != tc.expectedResp.Result {
				t.Errorf("Expected result %d, got %d", tc.expectedResp.Result, response.Result)
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[0:len(substr)] == substr
}
