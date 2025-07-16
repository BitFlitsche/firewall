#!/bin/bash
set -e

# Step 1: Build React frontend
cd firewall-app

echo "[1/3] Installing frontend dependencies..."
npm install

echo "[2/3] Building React frontend..."
npm run build

cd ..

# Step 2: Build Go backend
BINARY_NAME="firewall"
if [[ "$(uname -s)" == MINGW* || "$(uname -s)" == CYGWIN* || "$(uname -s)" == MSYS* ]]; then
  BINARY_NAME="firewall.exe"
fi

echo "[3/3] Building Go backend..."
go build -o "$BINARY_NAME" main.go

echo "Build complete. Starting the Go server..."

# Step 3: Start the Go binary
./$BINARY_NAME 