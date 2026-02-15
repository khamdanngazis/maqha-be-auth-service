# gRPC Service Testing Guide

## Overview
The auth service exposes a gRPC service on port 50053 with the following RPC:

- `GetUser(token: string)` - Validates a token and returns user information

## Testing Methods

### Method 1: Using grpcurl (Recommended)

#### Installation
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

#### List Available Services
```bash
grpcurl -plaintext localhost:50053 list
```

#### Test GetUser RPC - Positive Case
```bash
grpcurl -plaintext \
  -d '{"token": "admin_token_12345"}' \
  localhost:50053 model.User.GetUser
```

#### Test GetUser RPC - Invalid Token
```bash
grpcurl -plaintext \
  -d '{"token": "invalid_token"}' \
  localhost:50053 model.User.GetUser
```

#### Test GetUser RPC - Against Railway
```bash
# Note: Railway only exposes HTTP port, so direct gRPC connection won't work
# Instead, test the gRPC service inside the container or locally
grpcurl -plaintext maqha-be-auth-service-production.up.railway.app:50053 list
```

### Method 2: Using Go Test Client

Create a Go client to test:

```go
package main

import (
	"context"
	pb "maqhaa/auth_service/internal/interface/grpc/model"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserClient(conn)
	
	// Test with valid token
	resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{
		Token: "admin_token_12345",
	})
	if err != nil {
		log.Fatalf("Error calling GetUser: %v", err)
	}
	
	log.Printf("Response: Code=%d, Message=%s, UserID=%d", 
		resp.Code, resp.Message, resp.Data.Id)
}
```

### Method 3: Using Docker with grpcurl

For Railway deployment testing from local Docker:

```bash
docker run --rm \
  -e GRPC_ENDPOINT=maqha-be-auth-service-production.up.railway.app:50053 \
  fullstorydev/grpcurl:latest \
  -plaintext $GRPC_ENDPOINT list
```

## Test Cases

### TC-1: GetUser with Valid Token
- **Input:** token = "admin_token_12345"
- **Expected Response:**
  ```json
  {
    "code": 0,
    "message": "Success",
    "data": {
      "id": 1,
      "client_id": 1,
      "is_admin": true,
      "is_login": true
    }
  }
  ```

### TC-2: GetUser with Invalid Token
- **Input:** token = "invalid_token"
- **Expected Response:**
  ```json
  {
    "code": 202,
    "message": "Invalid Token",
    "data": null
  }
  ```

### TC-3: GetUser with Empty Token
- **Input:** token = ""
- **Expected Response:**
  ```json
  {
    "code": 202,
    "message": "Invalid Token",
    "data": null
  }
  ```

## Running Tests Locally

### 1. Start the service locally
```bash
go run cmd/main.go
```

### 2. In another terminal, run grpcurl tests
```bash
# List services
grpcurl -plaintext localhost:50053 list

# Test GetUser - valid token
grpcurl -plaintext \
  -d '{"token": "admin_token_12345"}' \
  localhost:50053 model.User.GetUser

# Test GetUser - invalid token
grpcurl -plaintext \
  -d '{"token": "invalid"}' \
  localhost:50053 model.User.GetUser
```

## Running Tests in Docker

### 1. Build the Docker image
```bash
docker build -t maqha-auth-service:latest .
```

### 2. Start the container
```bash
docker run -d \
  --name auth-service \
  -p 8011:8011 \
  -p 50053:50053 \
  -e AUTH_LOG_TO_STDOUT=true \
  -e AUTH_DATABASE_HOST=host.docker.internal \
  maqha-auth-service:latest
```

### 3. Test gRPC from host
```bash
grpcurl -plaintext localhost:50053 list
```

### 4. Stop the container
```bash
docker stop auth-service
docker rm auth-service
```

## Expected gRPC Service Structure

```protobuf
service User {
  rpc GetUser (GetUserRequest) returns (GetUserResponse);
}

message GetUserRequest {
  string token = 1;
}

message UserData {
  uint32 id = 1;
  uint32 client_id = 2;
  bool is_admin = 3;
  bool is_login = 4;
}

message GetUserResponse {
  int32 code = 1;
  string message = 2;
  UserData data = 3;
}
```

## Troubleshooting

### Connection Refused
- Ensure the service is running on port 50053
- Check if firewall is blocking the port
- Verify gRPC server is initialized correctly

### Service Not Listed
- The service might require reflection enabled
- Check if proto files are correctly compiled
- Ensure gRPC server is properly registered

### Invalid Token Response
- Verify the token exists in the database
- Check if token validation logic is working
- Ensure database connection is established

## Notes
- Railway platform only exposes HTTP ports publicly (port 8011)
- gRPC port (50053) is only accessible from within the Railway environment
- For production testing against Railway, use an internal network or deploy gRPC client in Railway
