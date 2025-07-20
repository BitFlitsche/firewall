#!/bin/bash

echo "Testing ASN Filter Implementation..."

# Test 1: Create an ASN rule
echo "Test 1: Creating ASN rule..."
curl -X POST http://localhost:8081/api/asn \
  -H "Content-Type: application/json" \
  -d '{
    "asn": "AS7922",
    "name": "Comcast Cable Communications",
    "status": "denied"
  }'

echo -e "\n\nTest 2: Getting ASN rules..."
curl -X GET "http://localhost:8081/api/asns?page=1&limit=10"

echo -e "\n\nTest 3: Testing ASN filtering with IP..."
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "68.85.108.1",
    "email": "",
    "user_agent": "",
    "country": "",
    "username": ""
  }'

echo -e "\n\nTest 4: Testing ASN filtering with different IP..."
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "email": "",
    "user_agent": "",
    "country": "",
    "username": ""
  }'

echo -e "\n\nTest 5: Testing ASN filtering with manual ASN override..."
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "8.8.8.8",
    "email": "",
    "user_agent": "",
    "country": "",
    "asn": "AS7922",
    "username": ""
  }'

echo -e "\n\nASN Filter Testing Complete!" 