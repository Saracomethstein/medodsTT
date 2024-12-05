#!/bin/bash

print_result() {
  if [ "$1" = "SUCCESS" ]; then
    echo -e "\033[32m[SUCCESS]\033[0m $2"
  else
    echo -e "\033[31m[FAILED]\033[0m $2"
  fi
}

# Test 1: Correct inputs
echo "Running Test 1: Correct inputs"
generate_token_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "2aab4e76-42ce-2de2-9e7a-a02e52e6eed8", "ip": "192.168.1.4"}')

echo "Response for generate-token (correct input):"
echo generate_token_response | jq .

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
echo "Running Test 2: Invalid inputs"
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
echo "Running Test 3: Invalid access token"
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
echo "Running Test 4: Empty user ID"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "", "ip": "192.168.1.4"}')

echo "Response for generate-token (empty user_id):"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for empty user_id."
else
  print_result "FAILED" "Unexpected behavior for empty user_id."
fi

# Test 5: Invalid IP address format
echo "Running Test 5: Invalid IP address format"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"user_id": "2aab4e76-42ce-2de2-9e7a-a02e52e6eed8", "ip": "999.999.999.999"}')

echo "Response for generate-token (invalid IP):"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for invalid IP address."
else
  print_result "FAILED" "Unexpected behavior for invalid IP address."
fi

# Test 6: Empty refresh token
echo "Running Test 6: Empty refresh token"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_token", "refresh_token": ""}')

echo "Response for validate-token (empty refresh_token):"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for empty refresh_token."
else
  print_result "FAILED" "Unexpected behavior for empty refresh_token."
fi

# Test 7: Missing access token
echo "Running Test 7: Missing access token"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "", "refresh_token": "valid_refresh_token"}')

echo "Response for missing access token:"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for missing access token."
else
  print_result "FAILED" "Unexpected behavior for missing access token."
fi

# Test 8: Invalid refresh token
echo "Running Test 8: Invalid refresh token"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_access_token", "refresh_token": "invalid_refresh_token"}')

echo "Response for invalid refresh token:"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for invalid refresh token."
else
  print_result "FAILED" "Unexpected behavior for invalid refresh token."
fi

# Test 9: Expired access token
echo "Running Test 9: Expired access token"
expired_access_token="expired.token.here"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$expired_access_token\", \"refresh_token\": \"$refresh_token\"}")

echo "Response for expired access token:"
echo $validate_token_error_response | jq .

if echo $validate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for expired access token."
else
  print_result "FAILED" "Unexpected behavior for expired access token."
fi

# Test 10: Access token with invalid structure
echo "Running Test 10: Access token with invalid structure"
invalid_structure_access_token="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.invalidstructure"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$invalid_structure_access_token\", \"refresh_token\": \"$refresh_token\"}")

echo "Response for malformed access token:"
echo $validate_token_error_response | jq .

if echo $validate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for malformed access token."
else
  print_result "FAILED" "Unexpected behavior for malformed access token."
fi

# Test 11: Valid refresh token, expired access token
echo "Running Test 11:  Valid refresh token, expired access token"
expired_access_token="expired.token.here"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$expired_access_token\", \"refresh_token\": \"$refresh_token\"}")

echo "Response for expired access token with valid refresh token:"
echo $validate_token_error_response | jq .

if echo $validate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for expired access token with valid refresh token."
else
  print_result "FAILED" "Unexpected behavior for expired access token with valid refresh token."
fi

# Test 12: Refresh with valid but non-matching tokens
echo "Running Test 12: Refresh with valid but non-matching tokens"
non_matching_refresh_token="non_matching_refresh_token"
validate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"access_token\": \"$access_token\", \"refresh_token\": \"$non_matching_refresh_token\"}")

echo "Response for non-matching refresh token:"
echo $validate_token_error_response | jq .

if echo $validate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for non-matching refresh token."
else
  print_result "FAILED" "Unexpected behavior for non-matching refresh token."
fi

# Test 13: Multiple refresh attempts with invalid tokens
echo "Running Test 13: Multiple refresh attempts with invalid tokens"
multiple_invalid_attempts_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "invalid_token", "refresh_token": "invalid_refresh_token"}')

echo "Response for multiple invalid refresh attempts (1):"
echo $multiple_invalid_attempts_response | jq .

if echo $multiple_invalid_attempts_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for the first invalid refresh attempt."
else
  print_result "FAILED" "Unexpected behavior for the first invalid refresh attempt."
fi

multiple_invalid_attempts_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "invalid_token", "refresh_token": "invalid_refresh_token"}')

echo "Response for multiple invalid refresh attempts (2):"
echo $multiple_invalid_attempts_response | jq .

if echo $multiple_invalid_attempts_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for the second invalid refresh attempt."
else
  print_result "FAILED" "Unexpected behavior for the second invalid refresh attempt."
fi

# Test 14: Missing refresh token in request
echo "Running Test 14: Missing refresh token in request"
generate_token_error_response=$(curl -s -X POST "http://localhost:8000/auth/refresh" \
    -H "Content-Type: application/json" \
    -d '{"access_token": "valid_access_token"}')

echo "Response for missing refresh token in request:"
echo $generate_token_error_response | jq .

if echo $generate_token_error_response | jq -e '.message' > /dev/null; then
  print_result "SUCCESS" "Error returned as expected for missing refresh token in request."
else
  print_result "FAILED" "Unexpected behavior for missing refresh token in request."
fi

