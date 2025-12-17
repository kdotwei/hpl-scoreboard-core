#!/bin/bash

# Test Pagination API Integration
# 这个脚本测试新的分页 API 功能

echo "Building the application..."
go build -o api ./cmd/api

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Starting server in background..."
./api &
SERVER_PID=$!

# 等待服务器启动
sleep 3

# 测试基本分页功能
echo "Testing pagination API..."
echo "1. Testing basic pagination (no cursor):"
curl -s "http://localhost:8080/api/v1/scores/paginated?limit=5" | jq '.'

echo -e "\n2. Testing with limit parameter:"
curl -s "http://localhost:8080/api/v1/scores/paginated?limit=2" | jq '.'

echo -e "\n3. Testing original endpoint (backward compatibility):"
curl -s "http://localhost:8080/api/v1/scores?limit=3" | jq '.'

# 清理
echo -e "\nStopping server..."
kill $SERVER_PID

echo "Integration test completed!"