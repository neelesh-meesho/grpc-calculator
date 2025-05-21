# === CONFIGURATION ===

# Proto directories
SERVER1_PROTO_DIR = server1/calc
SERVER2_PROTO_DIR = server2/calc

# Proto files
SERVER1_PROTO = $(SERVER1_PROTO_DIR)/calc.proto
SERVER2_PROTO = $(SERVER2_PROTO_DIR)/calc.proto

# Output binary directory
BIN_DIR = bin

# Protoc plugins
PROTOC_GEN_GO = /Users/neeleshpandey/go/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC = /Users/neeleshpandey/go/bin/protoc-gen-go-grpc

# Server port numbers
SERVER1_PORT = 50051
SERVER2_PORT = 50052

# === TARGETS ===

# Kill all servers
kill-all-servers:
	@echo "🛑 Stopping all running servers..."
	@-pkill -f bin/server || true

# Kill server1 if running
kill-server1:
	@echo "🛑 Stopping server1 if running..."
	@-pkill -f bin/server1 || true

# Kill server2 if running
kill-server2:
	@echo "🛑 Stopping server2 if running..."
	@-pkill -f bin/server2 || true

# Generate gRPC code for server1
generate-server1:
	@echo "🔧 Generating gRPC code for server1..."
	PATH="$$PATH:/Users/neeleshpandey/go/bin" protoc \
		--proto_path=$(SERVER1_PROTO_DIR) \
		--go_out=$(SERVER1_PROTO_DIR) \
		--go-grpc_out=$(SERVER1_PROTO_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(SERVER1_PROTO)

# Generate gRPC code for server2
generate-server2:
	@echo "🔧 Generating gRPC code for server2..."
	PATH="$$PATH:/Users/neeleshpandey/go/bin" protoc \
		--proto_path=$(SERVER2_PROTO_DIR) \
		--go_out=$(SERVER2_PROTO_DIR) \
		--go-grpc_out=$(SERVER2_PROTO_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(SERVER2_PROTO)

# Generate for both
generate: generate-server1 generate-server2

# Build server1
build-server1: generate-server1
	@echo "📦 Building server1..."
	cd server1 && go build -o ../$(BIN_DIR)/server1 .

# Build server2
build-server2: generate-server2
	@echo "📦 Building server2..."
	cd server2 && go build -o ../$(BIN_DIR)/server2 .

# Build both
build: build-server1 build-server2

# Run server1
run-server1: build-server1 kill-server1
	@echo "🚀 Running server1..."
	./$(BIN_DIR)/server1

# Run server2
run-server2: build-server2 kill-server2
	@echo "🚀 Running server2..."
	./$(BIN_DIR)/server2

# Run both servers in background
run-servers: build kill-all-servers
	@echo "🚀 Starting both servers in background..."
	@./$(BIN_DIR)/server1 &
	@./$(BIN_DIR)/server2 &
	@echo "✅ Both servers are running."

# Check if servers are running
check-servers:
	@echo "📊 Checking server status..."
	@echo "Server1 (port $(SERVER1_PORT)):"
	@-lsof -i:$(SERVER1_PORT) || echo "Not running"
	@echo "Server2 (port $(SERVER2_PORT)):"
	@-lsof -i:$(SERVER2_PORT) || echo "Not running"

# Clean generated files and binaries
clean: kill-all-servers
	@echo "🧹 Cleaning generated files and binaries..."
	rm -rf $(BIN_DIR)
	find server1/calc -name '*.pb.go' -delete
	find server2/calc -name '*.pb.go' -delete

# Create required directories
init:
	@echo "🛠️ Creating necessary directories..."
	mkdir -p $(BIN_DIR)

# === TESTING TARGETS ===

# Run tests for server1
test-server1: generate-server1
	@echo "🧪 Running tests for server1..."
	cd server1 && go test -v ./...

# Run tests for server2
test-server2: generate-server2
	@echo "🧪 Running tests for server2..."
	cd server2 && go test -v ./...

# Run all tests
test: test-server1 test-server2
	@echo "✅ All tests completed"

# Run tests with coverage
test-coverage: generate
	@echo "🧪 Running tests with coverage..."
	cd server1 && go test -v -coverprofile=../coverage-server1.out ./...
	cd server2 && go test -v -coverprofile=../coverage-server2.out ./...
	echo "mode: set" > coverage.out
	tail -n +2 coverage-server1.out >> coverage.out
	tail -n +2 coverage-server2.out >> coverage.out
	rm coverage-server1.out coverage-server2.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated in coverage.html"

# === ONE-COMMAND TO DO IT ALL ===

all: init generate build
	@echo "🎉 All steps completed!"
	@echo "📋 To run the servers:"
	@echo "  - Single terminal: make run-servers"
	@echo "  - Separate terminals: make run-server1 and make run-server2"

.PHONY: generate generate-server1 generate-server2 build build-server1 build-server2 \
	run-server1 run-server2 run-servers clean init all \
	kill-all-servers kill-server1 kill-server2 check-servers \
	test test-server1 test-server2 test-coverage
