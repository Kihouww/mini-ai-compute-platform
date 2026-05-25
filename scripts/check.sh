#!/usr/bin/env bash

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_KEY="${API_KEY:-test-api-key}"

echo "== health check =="
curl -s "$BASE_URL/health"
echo
echo

echo "== chat request =="
curl -s -X POST "$BASE_URL/v1/chat" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $API_KEY" \
  -d '{"model":"mock-llm","prompt":"check script"}'
echo
echo

echo "== request logs =="
curl -s "$BASE_URL/v1/requests" \
  -H "Authorization: Bearer $API_KEY"
echo
echo

echo "== rate limit smoke test =="
for i in $(seq 1 25); do
  code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/v1/chat" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $API_KEY" \
    -d "{\"model\":\"mock-llm\",\"prompt\":\"rate limit check $i\"}")
  echo "request $i -> $code"
done
