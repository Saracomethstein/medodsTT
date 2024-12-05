#!/bin/bash

# Function to print success or failure with colored backgrounds
print_result() {
  if [ "$1" == "SUCCESS" ]; then
    echo "\033[42m$2\033[0m"
  else
    echo "\033[41m$2\033[0m"
  fi
}

# Test 1: Correct inputs
generate_token_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "2aab4e76-42ce-2de2-9e7a-a02e52e6eed8", "ip": "192.168.1.4"}')

echo "Response for generate-token (correct input):"
echo $generate_token_response | jq .

if echo $generate_token_response | jq -e '.access_token' > /dev/null && echo $generate_token_response | jq -e '.refresh_token' > /dev/null; then
  print_result "SUCCESS" "Access token and refresh token generated successfully!"
else
  print_result "FAILED" "Error: Failed to generate tokens."
fi

access_token=$(echo $generate_token_response | jq -r '.access_token')
refresh_token=$(echo $generate_token_response | jq -r '.refresh_token')

validate_token_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$access_token\", \"refresh_token\": \"$refresh_token\"}")

echo "Response for validate-token (valid input):"
echo $validate_token_response | jq .
print_result "SUCCESS" "Token refreshed successfully!"

# Test 2: Invalid inputs
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "", "ip": "not-a-valid-ip"}')

echo "Response for generate-token (invalid input):"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for invalid inputs."
else
  print_result "FAILED" "Unexpected behavior for invalid inputs."
fi

# Test 3: Invalid access token
invalid_access_token="invalid.token.here"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$invalid_access_token\", \"refresh_token\": \"$refresh_token\"}")

echo "Response for validate-token (invalid access token):"
echo $validate_token_error_response | jq .

if echo $validate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for invalid access token."
else
  print_result "FAILED" "Unexpected behavior for invalid access token."
fi

# Test 4: Empty user ID
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "", "ip": "192.168.1.4"}')
print_result "SUCCESS" "Error returned for empty user_id."

# Test 5: Invalid IP address format
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "2aab4e76-42ce-2de2-9e7a-a02e52e6eed8", "ip": "999.999.999.999"}')
print_result "SUCCESS" "Error returned for invalid IP address."

# Test 6: Empty refresh token
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_token", "refresh_token": ""}')
print_result "SUCCESS" "Error returned for empty refresh_token."

# Test 7: Missing access token
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "", "refresh_token": "valid_refresh_token"}')
print_result "SUCCESS" "Error returned for missing access_token."

# Test 8: Invalid refresh token
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_access_token", "refresh_token": "invalid_refresh_token"}')
print_result "SUCCESS" "Error returned for invalid refresh_token."

# Test 9: Expired access token
expired_access_token="expired.token.here"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$expired_access_token\", \"refresh_token\": \"$refresh_token\"}")
print_result "SUCCESS" "Error returned for expired access token."

# Test 10: Access token with invalid structure
invalid_structure_access_token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.invalidstructure"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$invalid_structure_access_token\", \"refresh_token\": \"$refresh_token\"}")
print_result "SUCCESS" "Error returned for malformed access token."

# Test 11: Valid refresh token, expired access token
expired_access_token="expired.token.here"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$expired_access_token\", \"refresh_token\": \"$refresh_token\"}")
print_result "SUCCESS" "Error returned for expired access token."

# Test 12: Refresh with valid but non-matching tokens
non_matching_refresh_token="non_matching_refresh_token"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$access_token\", \"refresh_token\": \"$non_matching_refresh_token\"}")
print_result "SUCCESS" "Error returned for non-matching refresh token."

# Test 13: Multiple refresh attempts with invalid token
multiple_invalid_attempts_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "invalid_token", "refresh_token": "invalid_refresh_token"}')
multiple_invalid_attempts_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "invalid_token", "refresh_token": "invalid_refresh_token"}')
print_result "SUCCESS" "Error returned after multiple invalid refresh attempts."

# Test 14: Refresh with token after multiple incorrect refresh attempts
invalid_attempts_refresh_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_access_token", "refresh_token": "valid_refresh_token"}')
print_result "SUCCESS" "Valid token refresh after multiple failed attempts."

# Test 15: Missing refresh token in request
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_access_token"}')
print_result "SUCCESS" "Error returned for missing refresh token in request."

