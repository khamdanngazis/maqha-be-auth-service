#!/bin/bash

# gRPC Service Test Script for Railway Deployment
# This script tests the gRPC service using grpcurl

set -e

echo "=== maqha-be-auth-service gRPC Testing ==="
echo ""

# Check if grpcurl is available
if ! command -v grpcurl &> /dev/null; then
    echo "Installing grpcurl..."
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
fi

GRPC_ENDPOINT="${1:-localhost:50053}"
echo "Testing gRPC endpoint: $GRPC_ENDPOINT"
echo ""

# Test 1: List available services
echo "=== Test 1: List Available Services ==="
grpcurl -plaintext "$GRPC_ENDPOINT" list || echo "Note: Service reflection might not be enabled"
echo ""

# Test 2: GetUser with valid token
echo "=== Test 2: GetUser with Valid Token (admin) ==="
grpcurl -plaintext \
  -d '{"token": "admin_token_12345"}' \
  "$GRPC_ENDPOINT" model.User.GetUser
echo ""

# Test 3: GetUser with another valid token
echo "=== Test 3: GetUser with Valid Token (staff) ==="
grpcurl -plaintext \
  -d '{"token": "staff_token_12345"}' \
  "$GRPC_ENDPOINT" model.User.GetUser
echo ""

# Test 4: GetUser with invalid token
echo "=== Test 4: GetUser with Invalid Token ==="
grpcurl -plaintext \
  -d '{"token": "invalid_token_xyz"}' \
  "$GRPC_ENDPOINT" model.User.GetUser
echo ""

# Test 5: GetUser with empty token
echo "=== Test 5: GetUser with Empty Token ==="
grpcurl -plaintext \
  -d '{"token": ""}' \
  "$GRPC_ENDPOINT" model.User.GetUser
echo ""

# Test 6: GetUser with expired token
echo "=== Test 6: GetUser with Expired Token ==="
grpcurl -plaintext \
  -d '{"token": "expiredadmin_token_12345"}' \
  "$GRPC_ENDPOINT" model.User.GetUser
echo ""

echo "=== gRPC Tests Completed ==="
