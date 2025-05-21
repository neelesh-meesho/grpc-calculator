# gRPC Calculator

This project demonstrates a simple gRPC setup with two servers:
- `server2` - Performs the actual calculations
- `server1` - Acts as an intermediary that forwards requests to `server2`

## Prerequisites

1. Install necessary tools:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

2. Make sure `protoc` is installed and your `$GOPATH/bin` is in your `PATH`:
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

## Building and Running

This project uses a Makefile to simplify the build and run process. Here are the main commands:

### Build the Project

```bash
# Initialize, generate code, and build both servers
make all
```

### Run the Servers

```bash
# Run both servers in the background (in a single terminal)
make run-servers

# Or run each server individually (in separate terminals)
make run-server1  # first terminal
make run-server2  # second terminal
```

### Clean Up

```bash
# Stop all running servers and clean generated files
make clean
```

### Other Useful Commands

```bash
# Generate protobuf code
make generate

# Build executables only
make build

# Stop any running servers
make kill-servers
```

## Testing

### With Postman

1. Open Postman and create a new gRPC request
2. Set the server URL to `localhost:50051`
3. In Method Selection:
   - Select "main.CalculateService/Calculate" from the service definition
   - If importing proto file is required, use the proto file from `server1/calc/calc.proto`
4. In Message section, use this JSON format:
   ```json
   {
     "a": 10,
     "opr": "+",
     "b": 5
   }
   ```
5. Click "Invoke" to send the request

### With grpcurl

You can also use tools like `grpcurl` to test the service:
```bash
grpcurl -plaintext -d '{"a": 10, "opr": "+", "b": 5}' localhost:50051 main.CalculateService/Calculate
```

## Project Structure

- `server1/calc/calc.proto` - Protocol buffer definition
- `server1/main.go` - Server1 implementation that forwards requests to server2
- `server2/calc/calc.proto` - Same protocol buffer definition
- `server2/main.go` - Server2 implementation that performs calculations
- `Makefile` - Build automation
- `bin/` - Output directory for binaries 