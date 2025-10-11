#!/bin/bash

# =============================================================================
# Deployment Verification Script
# Run this after deployment to verify everything works
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_test() {
    echo -e "${BLUE}Testing:${NC} $1"
}

print_pass() {
    echo -e "${GREEN}✅ PASS:${NC} $1"
}

print_fail() {
    echo -e "${RED}❌ FAIL:${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}⚠️  WARN:${NC} $1"
}

FAILED_TESTS=0
PASSED_TESTS=0

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          Deployment Verification Tests                     ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Test 1: Check if Go API is running on port 8080
print_test "Go API on localhost:8080"
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080 | grep -E "200|404|401" > /dev/null; then
    print_pass "Go API is responding"
    ((PASSED_TESTS++))
else
    print_fail "Go API is not responding"
    ((FAILED_TESTS++))
fi

# Test 2: Check if existing Node app is running on port 2299
print_test "Existing Node.js app on localhost:2299"
if curl -s -o /dev/null -w "%{http_code}" http://localhost:2299 | grep -E "200|404|401" > /dev/null; then
    print_pass "Existing app is still running"
    ((PASSED_TESTS++))
else
    print_fail "Existing app is not responding"
    ((FAILED_TESTS++))
fi

# Test 3: Check PM2 processes
print_test "PM2 process management"
if pm2 list | grep -q "nutriediet-go-api"; then
    print_pass "New Go API is in PM2"
    ((PASSED_TESTS++))
else
    print_fail "Go API not found in PM2"
    ((FAILED_TESTS++))
fi

if pm2 list | grep -q "app"; then
    print_pass "Existing app is still in PM2"
    ((PASSED_TESTS++))
else
    print_warn "Existing 'app' not found in PM2 (may have different name)"
fi

# Test 4: Check if directories exist
print_test "Directory structure"
if [ -d "/home/sk/mys/nutriediet-new/backend" ]; then
    print_pass "Backend directory exists"
    ((PASSED_TESTS++))
else
    print_fail "Backend directory not found"
    ((FAILED_TESTS++))
fi

if [ -d "/home/sk/mys/nutriediet-new/frontend/build" ]; then
    print_pass "Frontend build directory exists"
    ((PASSED_TESTS++))
else
    print_fail "Frontend build directory not found"
    ((FAILED_TESTS++))
fi

# Test 5: Check if binary exists and is executable
print_test "Go binary"
if [ -x "/home/sk/mys/nutriediet-new/backend/nutriediet-go" ]; then
    print_pass "Go binary is executable"
    ((PASSED_TESTS++))
else
    print_fail "Go binary not found or not executable"
    ((FAILED_TESTS++))
fi

# Test 6: Check if .env file exists
print_test "Backend configuration"
if [ -f "/home/sk/mys/nutriediet-new/backend/.env" ]; then
    print_pass ".env file exists"
    ((PASSED_TESTS++))
else
    print_fail ".env file not found"
    ((FAILED_TESTS++))
fi

# Test 7: Check Nginx configuration
print_test "Nginx configuration"
if sudo nginx -t 2>&1 | grep -q "successful"; then
    print_pass "Nginx configuration is valid"
    ((PASSED_TESTS++))
else
    print_fail "Nginx configuration has errors"
    ((FAILED_TESTS++))
fi

# Test 8: Check if Nginx is running
print_test "Nginx service"
if systemctl is-active --quiet nginx; then
    print_pass "Nginx is running"
    ((PASSED_TESTS++))
else
    print_fail "Nginx is not running"
    ((FAILED_TESTS++))
fi

# Test 9: Check database connection (if possible)
print_test "Database connectivity"
if mysql -u nutriediet_new_user -p"${DB_PASSWORD}" -e "USE nutriediet_new_db;" 2>/dev/null; then
    print_pass "Database connection successful"
    ((PASSED_TESTS++))
else
    print_warn "Cannot verify database (password may be needed)"
fi

# Test 10: Check log files
print_test "Log files"
if [ -d "/home/sk/mys/nutriediet-new/logs" ]; then
    print_pass "Logs directory exists"
    ((PASSED_TESTS++))
else
    print_warn "Logs directory not found"
fi

# Test 11: Check React build files
print_test "React build files"
if [ -f "/home/sk/mys/nutriediet-new/frontend/build/index.html" ]; then
    print_pass "React index.html exists"
    ((PASSED_TESTS++))
else
    print_fail "React index.html not found"
    ((FAILED_TESTS++))
fi

if [ -d "/home/sk/mys/nutriediet-new/frontend/build/static" ]; then
    print_pass "React static assets exist"
    ((PASSED_TESTS++))
else
    print_fail "React static directory not found"
    ((FAILED_TESTS++))
fi

# Summary
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 All critical tests passed!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Test in browser: https://nutriediet.com/new"
    echo "  2. Verify API: https://nutriediet.com/new/api"
    echo "  3. Check existing site: https://nutriediet.com"
    echo ""
    echo "Useful commands:"
    echo "  - View logs: pm2 logs nutriediet-go-api"
    echo "  - Check status: pm2 list"
    echo "  - Monitor: pm2 monit"
    exit 0
else
    echo -e "${RED}⚠️  Some tests failed. Please review the errors above.${NC}"
    echo ""
    echo "Troubleshooting:"
    echo "  - Check logs: pm2 logs nutriediet-go-api"
    echo "  - Check PM2: pm2 list"
    echo "  - Check Nginx: sudo nginx -t"
    echo "  - Check Nginx logs: sudo tail -f /var/log/nginx/error.log"
    exit 1
fi

