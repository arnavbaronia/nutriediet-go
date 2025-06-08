#!/bin/bash

# Password Reset Testing Script
# This script helps you test the forgot password functionality

BASE_URL="http://localhost:8080"
TEST_EMAIL="test@example.com"

echo "🧪 Password Reset Testing Tool"
echo "=============================="

# Function to test forgot password
test_forgot_password() {
    echo ""
    echo "📧 Testing Forgot Password Endpoint..."
    echo "Endpoint: POST $BASE_URL/auth/forgot-password"
    
    read -p "Enter email to test with [$TEST_EMAIL]: " email
    email=${email:-$TEST_EMAIL}
    
    echo "Sending request..."
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST "$BASE_URL/auth/forgot-password" \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$email\"}")
    
    body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')
    status=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    echo "Status Code: $status"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "Response: $body"
    
    if [ "$status" = "200" ]; then
        echo "✅ Success! Check your email for the OTP"
        return 0
    else
        echo "❌ Error occurred"
        return 1
    fi
}

# Function to test reset password
test_reset_password() {
    echo ""
    echo "🔐 Testing Reset Password Endpoint..."
    echo "Endpoint: POST $BASE_URL/auth/reset-password"
    
    read -p "Enter email: " email
    read -p "Enter OTP from email: " otp
    read -s -p "Enter new password: " new_password
    echo ""
    
    echo "Sending request..."
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST "$BASE_URL/auth/reset-password" \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$email\", \"otp\": \"$otp\", \"new_password\": \"$new_password\"}")
    
    body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')
    status=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    echo "Status Code: $status"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "Response: $body"
    
    if [ "$status" = "200" ]; then
        echo "✅ Password reset successful!"
        return 0
    else
        echo "❌ Password reset failed"
        return 1
    fi
}

# Function to test login with new password
test_login() {
    echo ""
    echo "🔑 Testing Login with New Password..."
    echo "Endpoint: POST $BASE_URL/login"
    
    read -p "Enter email: " email
    read -s -p "Enter password: " password
    echo ""
    read -p "Enter user type [CLIENT]: " user_type
    user_type=${user_type:-CLIENT}
    
    echo "Sending request..."
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d "{\"email\": \"$email\", \"password\": \"$password\", \"user_type\": \"$user_type\"}")
        
    body=$(echo $response | sed -E 's/HTTPSTATUS:[0-9]{3}$//')
    status=$(echo $response | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    echo "Status Code: $status"
    echo "Response: $body" | jq '.' 2>/dev/null || echo "Response: $body"
    
    if [ "$status" = "200" ]; then
        echo "✅ Login successful!"
        return 0
    else
        echo "❌ Login failed"
        return 1
    fi
}

# Main menu
while true; do
    echo ""
    echo "Choose an option:"
    echo "1. Test Forgot Password (Request OTP)"
    echo "2. Test Reset Password (Verify OTP & Reset)"
    echo "3. Test Login with New Password"
    echo "4. Full Test Flow"
    echo "5. Exit"
    
    read -p "Enter choice [1-5]: " choice
    
    case $choice in
        1)
            test_forgot_password
            ;;
        2)
            test_reset_password
            ;;
        3)
            test_login
            ;;
        4)
            echo "🚀 Running Full Test Flow..."
            if test_forgot_password; then
                echo ""
                echo "📧 OTP sent! Now check your email and continue with step 2..."
                read -p "Press Enter when you have the OTP..."
                test_reset_password
                echo ""
                echo "🔐 Now test login with your new password..."
                test_login
            fi
            ;;
        5)
            echo "👋 Goodbye!"
            exit 0
            ;;
        *)
            echo "❌ Invalid choice. Please try again."
            ;;
    esac
done 