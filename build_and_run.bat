@echo off
setlocal enabledelayedexpansion

REM Step 1: Build React frontend
cd firewall-app

echo [1/3] Installing frontend dependencies...
npm install || (echo npm install failed & exit /b 1)

echo [2/3] Building React frontend...
npm run build || (echo npm run build failed & exit /b 1)

cd ..

REM Step 2: Build Go backend
echo [3/3] Building Go backend...
go build -o firewall.exe main.go || (echo Go build failed & exit /b 1)

echo Build complete. Starting the Go server...

REM Step 3: Start the Go binary
firewall.exe 