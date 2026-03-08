#!/bin/bash
#
# Comprehensive API Test Suite for NutrieDiet Go Backend
# Simulates a real user journey through all endpoints.
#
# Rate-limit aware: the server has 5 req/min on auth routes and
# 3 req/min on password-reset routes, so we pace accordingly.
#
# Usage: bash test_api.sh
#

BASE_URL="http://localhost:8080/api"
PASS=0
FAIL=0
SKIP=0
TOKEN=""
REFRESH_TOKEN=""
CLIENT_ID=""
TIMESTAMP=$(date +%s)
TEST_EMAIL="testuser_${TIMESTAMP}@example.com"
TEST_PASSWORD="TestStr0ng!Pass99"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

print_header() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_sub() {
    echo ""
    echo -e "${YELLOW}--- $1 ---${NC}"
}

assert_status() {
    local test_name="$1" expected="$2" actual="$3" body="$4"
    if [ "$actual" -eq "$expected" ] 2>/dev/null; then
        echo -e "  ${GREEN}PASS${NC} $test_name (HTTP $actual)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC} $test_name (expected HTTP $expected, got HTTP $actual)"
        [ -n "$body" ] && echo -e "       Response: $(echo "$body" | cut -c1-200)"
        FAIL=$((FAIL + 1))
    fi
}

assert_oneof() {
    local test_name="$1" e1="$2" e2="$3" actual="$4" body="$5"
    if [ "$actual" -eq "$e1" ] 2>/dev/null || [ "$actual" -eq "$e2" ] 2>/dev/null; then
        echo -e "  ${GREEN}PASS${NC} $test_name (HTTP $actual)"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC} $test_name (expected HTTP $e1 or $e2, got HTTP $actual)"
        [ -n "$body" ] && echo -e "       Response: $(echo "$body" | cut -c1-200)"
        FAIL=$((FAIL + 1))
    fi
}

assert_contains() {
    local test_name="$1" text="$2" body="$3"
    if echo "$body" | grep -q "$text"; then
        echo -e "  ${GREEN}PASS${NC} $test_name"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC} $test_name (body missing '$text')"
        FAIL=$((FAIL + 1))
    fi
}

skip_test() {
    echo -e "  ${YELLOW}SKIP${NC} $1 — $2"
    SKIP=$((SKIP + 1))
}

wait_for_rate_limit() {
    local secs=${1:-62}
    echo -e "\n  ${YELLOW}Waiting ${secs}s for rate-limit window to reset...${NC}"
    sleep "$secs"
}

# Perform curl, set BODY and STATUS globals
do_req() {
    local method="$1" url="$2" data="$3" token="$4"
    local args=(-s -w "\n%{http_code}" -X "$method" "$url" -H "Content-Type: application/json")
    [ -n "$token" ] && args+=(-H "Authorization: Bearer $token")
    [ -n "$data" ] && args+=(-d "$data")
    local resp
    resp=$(curl "${args[@]}" 2>/dev/null)
    STATUS=$(echo "$resp" | tail -n 1)
    BODY=$(echo "$resp" | sed '$d')
}

echo -e "${BOLD}${CYAN}"
echo "  ╔═══════════════════════════════════════════════════════════╗"
echo "  ║          NutrieDiet API — Comprehensive Test Suite       ║"
echo "  ╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"
echo -e "  Base URL:   ${BASE_URL}"
echo -e "  Test email: ${TEST_EMAIL}"

########################################################################
# 1. SERVER CONNECTIVITY
########################################################################
print_header "1. SERVER CONNECTIVITY"

do_req "GET" "$BASE_URL/nonexistent"
if [ -n "$STATUS" ] && [ "$STATUS" -gt 0 ] 2>/dev/null; then
    echo -e "  ${GREEN}PASS${NC} Server is reachable"
    PASS=$((PASS + 1))
else
    echo -e "  ${RED}FAIL${NC} Server unreachable — start with: go run main.go"
    exit 1
fi

########################################################################
# 2. SECURITY HEADERS & CORS (no rate-limited endpoints)
########################################################################
print_header "2. SECURITY HEADERS & CORS"

print_sub "2a. Security Headers"
HEADERS=$(curl -s -D - -o /dev/null "$BASE_URL/nonexistent" 2>/dev/null)

for h in "X-Content-Type-Options" "X-Frame-Options" "X-XSS-Protection" "Content-Security-Policy"; do
    if echo "$HEADERS" | grep -qi "$h"; then
        echo -e "  ${GREEN}PASS${NC} $h header present"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}FAIL${NC} $h header missing"
        FAIL=$((FAIL + 1))
    fi
done

print_sub "2b. CORS — allowed origin"
CORS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X OPTIONS "$BASE_URL/login" \
    -H "Origin: https://nutriediet.netlify.app" \
    -H "Access-Control-Request-Method: POST" \
    -H "Access-Control-Request-Headers: Content-Type,Authorization")
assert_oneof "CORS preflight from allowed origin" 200 204 "$CORS_STATUS"

print_sub "2c. CORS — disallowed origin"
CORS_RESP=$(curl -s -D - -o /dev/null -X OPTIONS "$BASE_URL/login" \
    -H "Origin: https://evil-site.com" \
    -H "Access-Control-Request-Method: POST" 2>/dev/null)
if echo "$CORS_RESP" | grep -qi "access-control-allow-origin: https://evil-site.com"; then
    echo -e "  ${RED}FAIL${NC} CORS allows evil-site.com"
    FAIL=$((FAIL + 1))
else
    echo -e "  ${GREEN}PASS${NC} CORS blocks disallowed origin"
    PASS=$((PASS + 1))
fi

########################################################################
# 3. PROTECTED ROUTES WITHOUT TOKEN (uses protected group, not auth rate limit)
########################################################################
print_header "3. PROTECTED ROUTES WITHOUT AUTH TOKEN"

for route in "/users" "/admin/clients" "/clients/1/diet" "/admin/exercises" "/admin/recipes" "/admin/motivation"; do
    do_req "GET" "$BASE_URL$route"
    assert_oneof "GET $route without token" 400 401 "$STATUS" "$BODY"
done

########################################################################
# 4. PROTECTED ROUTES WITH INVALID TOKEN
########################################################################
print_header "4. PROTECTED ROUTES WITH INVALID TOKEN"

FAKE_TOKEN="eyJhbGciOiJIUzI1NiJ9.eyJlbWFpbCI6ImZha2VAZXhhbXBsZS5jb20iLCJleHAiOjE2MDAwMDAwMDB9.badsig"

do_req "GET" "$BASE_URL/users" "" "$FAKE_TOKEN"
assert_status "GET /users with invalid token" 401 "$STATUS" "$BODY"

do_req "GET" "$BASE_URL/admin/clients" "" "$FAKE_TOKEN"
assert_status "GET /admin/clients with invalid token" 401 "$STATUS" "$BODY"

########################################################################
# 5. SIGNUP + LOGIN (rate-limit budget: 5/min for auth endpoints)
#    Strategy: signup(1) + login(1) = 2 requests. Get the token first.
########################################################################
print_header "5. SIGNUP & LOGIN (rate-limit aware: 2 of 5 budget)"

print_sub "5a. Create user (valid)"
do_req "POST" "$BASE_URL/create_user" \
    "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"first_name\":\"Test\",\"last_name\":\"User\",\"user_type\":\"CLIENT\"}"
assert_oneof "Signup with valid data" 200 201 "$STATUS" "$BODY"

print_sub "5b. Login with valid credentials"
do_req "POST" "$BASE_URL/login" \
    "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"user_type\":\"CLIENT\"}"
assert_status "Login with valid credentials" 200 "$STATUS" "$BODY"

if [ "$STATUS" -eq 200 ]; then
    TOKEN=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('token',''))" 2>/dev/null)
    REFRESH_TOKEN=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('refreshToken',''))" 2>/dev/null)
    CLIENT_ID=$(echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('client_id',0))" 2>/dev/null)
    assert_contains "Response contains token" "token" "$BODY"
    assert_contains "Response contains refreshToken" "refreshToken" "$BODY"
    assert_contains "Response contains user_type" "user_type" "$BODY"
    assert_contains "Response contains email" "email" "$BODY"
    assert_contains "Response contains is_active" "is_active" "$BODY"
    echo -e "       client_id=$CLIENT_ID, first_time_login=$(echo "$BODY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('first_time_login','?'))" 2>/dev/null)"
fi

########################################################################
# 6. AUTHENTICATED CLIENT ROUTES (uses API rate limit: 100/min)
########################################################################
print_header "6. AUTHENTICATED CLIENT ROUTES"

if [ -z "$TOKEN" ]; then
    skip_test "All authenticated routes" "Login failed, no token"
else
    print_sub "6a. Cross-client access (should be blocked)"
    do_req "GET" "$BASE_URL/clients/999999/diet" "" "$TOKEN"
    assert_status "Access another client's diet" 401 "$STATUS" "$BODY"

    do_req "GET" "$BASE_URL/clients/999999/my_profile" "" "$TOKEN"
    assert_status "Access another client's profile" 401 "$STATUS" "$BODY"

    do_req "GET" "$BASE_URL/clients/999999/exercise" "" "$TOKEN"
    assert_status "Access another client's exercises" 401 "$STATUS" "$BODY"

    if [ -n "$CLIENT_ID" ] && [ "$CLIENT_ID" != "0" ]; then
        print_sub "6b. Own client data (client_id=$CLIENT_ID)"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/profile_created" "" "$TOKEN"
        assert_oneof "Check profile_created" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/my_profile" "" "$TOKEN"
        assert_oneof "Get own profile" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/diet" "" "$TOKEN"
        assert_oneof "Get own diet" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/exercise" "" "$TOKEN"
        assert_oneof "Get own exercises" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/weight-history" "" "$TOKEN"
        assert_oneof "Get weight history" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/recipe" "" "$TOKEN"
        assert_oneof "Get recipes" 200 404 "$STATUS" "$BODY"

        do_req "GET" "$BASE_URL/clients/$CLIENT_ID/motivation" "" "$TOKEN"
        assert_oneof "Get motivations" 200 404 "$STATUS" "$BODY"
    else
        skip_test "Own client data routes" "client_id=0 (first-time login, no client profile yet)"
    fi

    print_sub "6c. Client accessing admin routes (should be blocked)"
    do_req "GET" "$BASE_URL/admin/clients" "" "$TOKEN"
    assert_status "Client cannot access /admin/clients" 401 "$STATUS" "$BODY"

    do_req "GET" "$BASE_URL/admin/users" "" "$TOKEN"
    assert_status "Client cannot access /admin/users" 401 "$STATUS" "$BODY"
fi

########################################################################
# 7. AUTH NEGATIVE TESTS (wait for rate limit reset first)
########################################################################
print_header "7. AUTH NEGATIVE TESTS"
wait_for_rate_limit 62

print_sub "7a. Signup edge cases"

do_req "POST" "$BASE_URL/signup" '{"email":"bad"}'
assert_status "Signup with missing fields" 400 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/create_user" \
    "{\"email\":\"weak_${TIMESTAMP}@example.com\",\"password\":\"short\",\"first_name\":\"X\",\"last_name\":\"Y\",\"user_type\":\"CLIENT\"}"
assert_status "Signup with weak password" 400 "$STATUS" "$BODY"

LONG_PASS=$(python3 -c "print('Aa1!' + 'x'*200)")
do_req "POST" "$BASE_URL/create_user" \
    "{\"email\":\"long_${TIMESTAMP}@example.com\",\"password\":\"$LONG_PASS\",\"first_name\":\"X\",\"last_name\":\"Y\",\"user_type\":\"CLIENT\"}"
assert_status "Signup with >128 char password" 400 "$STATUS" "$BODY"

print_sub "7b. Login negative cases"
wait_for_rate_limit 62

do_req "POST" "$BASE_URL/login" \
    "{\"email\":\"$TEST_EMAIL\",\"password\":\"WrongPassword123!\",\"user_type\":\"CLIENT\"}"
assert_status "Login with wrong password" 403 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/login" \
    '{"email":"nonexistent@example.com","password":"whatever","user_type":"CLIENT"}'
assert_status "Login with non-existent email" 404 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/login" ''
assert_status "Login with empty body" 400 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/login" \
    "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"user_type\":\"ADMIN\"}"
assert_status "Login with wrong user_type" 404 "$STATUS" "$BODY"

########################################################################
# 8. FORGOT PASSWORD FLOW (3/min strict rate limit)
########################################################################
print_header "8. FORGOT PASSWORD FLOW"
wait_for_rate_limit 62

print_sub "8a. Request OTP — validation"
do_req "POST" "$BASE_URL/auth/forgot-password" '{"email":"not-an-email"}'
assert_status "Forgot password — invalid email format" 400 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/auth/forgot-password" '{}'
assert_status "Forgot password — empty body" 400 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/auth/forgot-password" '{"email":"ghost@example.com"}'
assert_status "Forgot password — non-existent email" 400 "$STATUS" "$BODY"

print_sub "8b. Request OTP — existing user (sends email)"
wait_for_rate_limit 62

do_req "POST" "$BASE_URL/auth/forgot-password" "{\"email\":\"$TEST_EMAIL\"}"
assert_oneof "Forgot password — existing user" 200 500 "$STATUS" "$BODY"
echo -e "       (500 is OK if SMTP not configured for test email domain)"

print_sub "8c. Reset password — validation"
wait_for_rate_limit 62

do_req "POST" "$BASE_URL/auth/reset-password" '{"email":"a@b.com"}'
assert_status "Reset — missing fields" 400 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/auth/reset-password" \
    "{\"email\":\"$TEST_EMAIL\",\"otp\":\"000000\",\"new_password\":\"NewStr0ng!Pass99\"}"
assert_oneof "Reset — wrong OTP" 400 500 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/auth/reset-password" \
    '{"email":"ghost@example.com","otp":"123456","new_password":"NewStr0ng!Pass99"}'
assert_status "Reset — non-existent email" 400 "$STATUS" "$BODY"

########################################################################
# 9. EDGE CASES & SECURITY
########################################################################
print_header "9. EDGE CASES & SECURITY"
wait_for_rate_limit 62

print_sub "9a. Wrong HTTP methods"
do_req "GET" "$BASE_URL/login"
assert_oneof "GET /login (should be POST)" 404 405 "$STATUS"

do_req "GET" "$BASE_URL/signup"
assert_oneof "GET /signup (should be POST)" 404 405 "$STATUS"

print_sub "9b. Injection attempts"
do_req "POST" "$BASE_URL/login" \
    '{"email":"admin@test.com OR 1=1","password":"test","user_type":"ADMIN"}'
assert_oneof "SQL injection in email" 400 404 "$STATUS" "$BODY"

do_req "POST" "$BASE_URL/login" \
    '{"email":"<script>alert(1)</script>@test.com","password":"test","user_type":"CLIENT"}'
assert_oneof "XSS in email" 400 404 "$STATUS" "$BODY"

print_sub "9c. Malformed requests"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" -d '{bad json}')
assert_status "Malformed JSON" 400 "$STATUS"

print_sub "9d. Malformed auth header"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/users" \
    -H "Authorization: notbearer sometoken")
assert_oneof "Auth header without Bearer prefix" 400 401 "$STATUS"

########################################################################
# 10. RATE LIMITING VERIFICATION
########################################################################
print_header "10. RATE LIMITING VERIFICATION"
wait_for_rate_limit 62

echo -e "  Sending 7 rapid requests to /api/login..."
RATE_LIMITED=false
for i in $(seq 1 7); do
    RL_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d '{"email":"rl@test.com","password":"x","user_type":"CLIENT"}')
    if [ "$RL_STATUS" -eq 429 ]; then
        RATE_LIMITED=true
        echo -e "  ${GREEN}PASS${NC} Rate limited after $i request(s) (HTTP 429)"
        PASS=$((PASS + 1))
        break
    fi
done
if [ "$RATE_LIMITED" = false ]; then
    echo -e "  ${RED}FAIL${NC} Rate limiter never triggered after 7 requests"
    FAIL=$((FAIL + 1))
fi

########################################################################
# SUMMARY
########################################################################
print_header "TEST RESULTS SUMMARY"

TOTAL=$((PASS + FAIL + SKIP))
echo ""
echo -e "  ${GREEN}Passed:  $PASS${NC}"
echo -e "  ${RED}Failed:  $FAIL${NC}"
echo -e "  ${YELLOW}Skipped: $SKIP${NC}"
echo -e "  ${BOLD}Total:   $TOTAL${NC}"
echo ""

if [ "$FAIL" -eq 0 ]; then
    echo -e "  ${GREEN}${BOLD}ALL TESTS PASSED${NC}"
elif [ "$FAIL" -le 3 ]; then
    echo -e "  ${YELLOW}${BOLD}MOSTLY PASSING — $FAIL minor issue(s) to review${NC}"
else
    echo -e "  ${RED}${BOLD}$FAIL TEST(S) FAILED — review above${NC}"
fi
echo ""
