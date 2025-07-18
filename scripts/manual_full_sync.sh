#!/bin/bash

# Manual Full Sync Script
# This script demonstrates how to trigger a manual full sync

echo "=== Manual Full Sync Demo ==="

# Check if the firewall service is running
if ! curl -s http://localhost:8081/system-stats > /dev/null; then
    echo "Error: Firewall service is not running on localhost:8081"
    echo "Please start the service first: ./firewall"
    exit 1
fi

echo "1. Current system stats:"
curl -s http://localhost:8081/system-stats | jq '.'

echo -e "\n2. Triggering manual full sync..."
response=$(curl -s -X POST http://localhost:8081/sync/full)

echo "Response: $response"

echo -e "\n3. Updated system stats:"
curl -s http://localhost:8081/system-stats | jq '.'

echo -e "\n=== Manual full sync completed ==="
echo "Note: This should only be used when:"
echo "- Initial setup"
echo "- Data recovery"
echo "- Schema changes"
echo "- Troubleshooting incremental sync issues" 