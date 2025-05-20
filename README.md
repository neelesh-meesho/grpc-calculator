# gRPC Calculator

This project demonstrates a simple gRPC setup with two servers:
- `server2` - Performs the actual calculations
- `server1` - Acts as an intermediary that forwards requests to `server2`

## Setup

1. Install necessary tools:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

2. Make sure `protoc` is installed and your `$GOPATH/bin` is in your `PATH`:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

## Running the application

1. Start server2 (calculation service) first:
   ```bash
   cd server2
   go run main.go
   ```

2. In a new terminal, start server1:
   ```bash
   cd server1
   go run main.go
   ```

3. Now you can send gRPC requests to server1 on port 50051, and it will forward them to server2 on port 50052.

## Testing

You can use a gRPC client like `grpcurl` or write a simple client application to test the service.

Example with grpcurl:
```bash
grpcurl -plaintext -d '{"a": 10, "opr": "+", "b": 5}' localhost:50051 main.CalculateService/Calculate
```

## Project Structure

- `server1/calc/calc.proto` - Protocol buffer definition
- `server1/main.go` - Server1 implementation that forwards requests to server2
- `server2/calc/calc.proto` - Same protocol buffer definition
- `server2/main.go` - Server2 implementation that performs calculations 