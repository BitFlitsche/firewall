#!/bin/bash

# Test script for Spamhaus ASN-DROP import functionality
# Make sure the firewall server is running on localhost:8081

echo "=== Spamhaus ASN-DROP Import Test ==="
echo "Note: Auto-import is disabled by default. Check config.yaml to enable."
echo

# Test 1: Get Spamhaus import stats (before import)
echo "Test 1: Getting Spamhaus import stats (before import)..."
curl -X GET http://localhost:8081/api/asns/spamhaus-stats
echo -e "\n"

# Test 2: Import Spamhaus ASN-DROP data
echo "Test 2: Importing Spamhaus ASN-DROP data..."
echo "This may take a moment as it fetches and processes data from Spamhaus..."
curl -X POST http://localhost:8081/api/asns/import-spamhaus
echo -e "\n"

# Test 3: Get Spamhaus import stats (after import)
echo "Test 3: Getting Spamhaus import stats (after import)..."
curl -X GET http://localhost:8081/api/asns/spamhaus-stats
echo -e "\n"

# Test 4: Check Spamhaus import status
echo "Test 4: Checking Spamhaus import status..."
curl -X GET http://localhost:8081/api/asns/spamhaus-status
echo -e "\n"



# Test 5: List ASNs with Spamhaus source
echo "Test 5: Listing ASNs with Spamhaus source..."
curl -X GET "http://localhost:8081/api/asns?page=1&limit=5&search=spamhaus"
echo -e "\n"

# Test 6: Test ASN filtering with imported Spamhaus data
echo "Test 6: Testing ASN filtering with imported Spamhaus data..."
curl -X POST http://localhost:8081/api/filter \
  -H "Content-Type: application/json" \
  -d '{
    "ip": "68.85.108.1",
    "email": "",
    "user_agent": "",
    "country": "",
    "asn": "",
    "username": ""
  }'
echo -e "\n"

# Test 7: Test manual ASN override with Spamhaus data
echo "Test 7: Testing manual ASN override with Spamhaus data..."
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
echo -e "\n"

echo "=== Spamhaus Import Testing Complete! ===" 