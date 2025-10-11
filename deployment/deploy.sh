#!/bin/bash

# =============================================================================
# Nutriediet New App Deployment Script
# Deploys Go backend + React frontend to www.nutriediet.com/new
# =============================================================================

set -e  # Exit on any error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DEPLOY_USER="sk"
DEPLOY_DIR="/home/sk/mys/nutriediet-new"
BACKEND_DIR="$DEPLOY_DIR/backend"
FRONTEND_DIR="$DEPLOY_DIR/frontend"
LOGS_DIR="$DEPLOY_DIR/logs"
GITHUB_REPO="https://github.com/cd-Ishita/nutriediet-go.git"
GITHUB_FRONTEND_REPO="YOUR_FRONTEND_REPO_URL"  # Update this
DB_NAME="nutriediet_new_db"
DB_USER="nutriediet_new_user"
GO_VERSION="1.21.5"
NODE_VERSION="20"

# Functions
print_step() {
    echo -e "${BLUE}==>${NC} ${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

check_user() {
    if [ "$USER" != "$DEPLOY_USER" ]; then
        print_error "This script must be run as user '$DEPLOY_USER'"
        exit 1
    fi
}

check_existing_site() {
    print_step "Checking existing site on port 2299..."
    if curl -s http://localhost:2299 > /dev/null; then
        print_success "Existing site is running on port 2299"
    else
        print_warning "Cannot reach existing site on port 2299. Is it running?"
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

install_go() {
    print_step "Checking Go installation..."
    
    if command -v go &> /dev/null; then
        CURRENT_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_success "Go $CURRENT_GO_VERSION is already installed"
    else
        print_step "Installing Go $GO_VERSION..."
        cd /tmp
        wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
        rm "go${GO_VERSION}.linux-amd64.tar.gz"
        
        # Add Go to PATH if not already there
        if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
            echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
            echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
        fi
        
        export PATH=$PATH:/usr/local/go/bin
        print_success "Go installed successfully"
    fi
}

upgrade_node() {
    print_step "Checking Node.js version..."
    CURRENT_NODE_VERSION=$(node -v | sed 's/v//' | cut -d'.' -f1)
    
    if [ "$CURRENT_NODE_VERSION" -lt 18 ]; then
        print_warning "Node.js version is $CURRENT_NODE_VERSION, upgrading to v$NODE_VERSION..."
        
        # Install nvm if not present
        if [ ! -d "$HOME/.nvm" ]; then
            curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
            export NVM_DIR="$HOME/.nvm"
            [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
        fi
        
        # Load nvm
        export NVM_DIR="$HOME/.nvm"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
        
        nvm install $NODE_VERSION
        nvm use $NODE_VERSION
        nvm alias default $NODE_VERSION
        print_success "Node.js upgraded to v$(node -v)"
    else
        print_success "Node.js v$(node -v) is already installed"
    fi
}

create_directories() {
    print_step "Creating deployment directories..."
    
    mkdir -p "$DEPLOY_DIR"
    mkdir -p "$BACKEND_DIR"
    mkdir -p "$FRONTEND_DIR"
    mkdir -p "$LOGS_DIR"
    mkdir -p "$BACKEND_DIR/images"
    
    print_success "Directories created"
}

setup_database() {
    print_step "Setting up MySQL database..."
    
    read -sp "Enter MySQL root password: " MYSQL_ROOT_PASSWORD
    echo
    
    read -sp "Enter new password for $DB_USER: " DB_PASSWORD
    echo
    
    # Create database and user
    mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON $DB_NAME.* TO '$DB_USER'@'localhost';
FLUSH PRIVILEGES;
EOF
    
    if [ $? -eq 0 ]; then
        print_success "Database created successfully"
        echo "$DB_PASSWORD" > /tmp/.db_password_temp
    else
        print_error "Failed to create database"
        exit 1
    fi
}

clone_repos() {
    print_step "Cloning repositories..."
    
    # Clone backend
    if [ -d "$BACKEND_DIR/.git" ]; then
        print_warning "Backend repo already exists, pulling latest changes..."
        cd "$BACKEND_DIR"
        git pull
    else
        print_step "Cloning backend repository..."
        git clone "$GITHUB_REPO" "$BACKEND_DIR"
    fi
    
    # Clone frontend
    if [ -d "$FRONTEND_DIR/.git" ]; then
        print_warning "Frontend repo already exists, pulling latest changes..."
        cd "$FRONTEND_DIR"
        git pull
    else
        print_step "Cloning frontend repository..."
        # Update this with your actual frontend repo
        echo "Note: Update GITHUB_FRONTEND_REPO in script with your frontend repository URL"
        # git clone "$GITHUB_FRONTEND_REPO" "$FRONTEND_DIR"
        # For now, copy local frontend
        print_warning "Copying local frontend files (update this in production)"
    fi
    
    print_success "Repositories cloned"
}

build_backend() {
    print_step "Building Go backend..."
    
    cd "$BACKEND_DIR"
    
    # Download dependencies
    go mod download
    go mod verify
    
    # Build binary
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nutriediet-go -ldflags="-s -w" .
    
    chmod +x nutriediet-go
    
    print_success "Backend built successfully"
}

configure_backend() {
    print_step "Configuring backend..."
    
    if [ -f /tmp/.db_password_temp ]; then
        DB_PASSWORD=$(cat /tmp/.db_password_temp)
        rm /tmp/.db_password_temp
    else
        read -sp "Enter database password: " DB_PASSWORD
        echo
    fi
    
    # Generate random secrets
    JWT_SECRET=$(openssl rand -hex 32)
    SESSION_SECRET=$(openssl rand -hex 32)
    
    # Create .env file
    cat > "$BACKEND_DIR/.env" <<EOF
# Server Configuration
PORT=8080
GIN_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME

# JWT Configuration
JWT_SECRET=$JWT_SECRET
JWT_EXPIRY=24h

# CORS Configuration
ALLOWED_ORIGINS=https://nutriediet.com,https://www.nutriediet.com

# File Upload
UPLOAD_DIR=/home/sk/mys/nutriediet-new/backend/images
MAX_UPLOAD_SIZE=10485760

# Application
APP_NAME=Nutriediet
APP_ENV=production
APP_URL=https://nutriediet.com/new
API_URL=https://nutriediet.com/new/api

# Security
SESSION_SECRET=$SESSION_SECRET
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
EOF
    
    chmod 600 "$BACKEND_DIR/.env"
    print_success "Backend configured"
}

build_frontend() {
    print_step "Building React frontend..."
    
    cd "$FRONTEND_DIR"
    
    # Create production env file
    cat > .env.production <<EOF
REACT_APP_API_URL=/new/api
PUBLIC_URL=/new
NODE_ENV=production
EOF
    
    # Add homepage to package.json if not present
    if ! grep -q '"homepage"' package.json; then
        # Backup original
        cp package.json package.json.backup
        # Add homepage field
        node -e "
        const pkg = require('./package.json');
        pkg.homepage = '/new';
        require('fs').writeFileSync('package.json', JSON.stringify(pkg, null, 2));
        "
    fi
    
    # Install dependencies
    npm ci --production=false
    
    # Build
    GENERATE_SOURCEMAP=false npm run build
    
    print_success "Frontend built successfully"
}

setup_pm2() {
    print_step "Setting up PM2..."
    
    # Install PM2 globally if not present
    if ! command -v pm2 &> /dev/null; then
        npm install -g pm2
    fi
    
    # Create ecosystem config
    cat > "$DEPLOY_DIR/ecosystem.config.js" <<EOF
module.exports = {
  apps: [
    {
      name: 'nutriediet-go-api',
      script: './nutriediet-go',
      cwd: '$BACKEND_DIR',
      instances: 1,
      exec_mode: 'fork',
      autorestart: true,
      watch: false,
      max_memory_restart: '500M',
      env: {
        PORT: '8080',
        GIN_MODE: 'release'
      },
      error_file: '$LOGS_DIR/go-api-error.log',
      out_file: '$LOGS_DIR/go-api-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss Z',
      merge_logs: true
    }
  ]
};
EOF
    
    # Stop if already running
    pm2 delete nutriediet-go-api 2>/dev/null || true
    
    # Start new app
    cd "$DEPLOY_DIR"
    pm2 start ecosystem.config.js
    
    # Save PM2 state
    pm2 save
    
    # Setup startup script if not done
    pm2 startup systemd -u $DEPLOY_USER --hp /home/$DEPLOY_USER
    
    print_success "PM2 configured and app started"
}

configure_nginx() {
    print_step "Updating Nginx configuration..."
    
    print_warning "Manual step required:"
    echo "1. Backup current nginx config:"
    echo "   sudo cp /etc/nginx/sites-available/nutriediet.com /etc/nginx/sites-available/nutriediet.com.backup"
    echo ""
    echo "2. The new nginx configuration is available at:"
    echo "   $BACKEND_DIR/deployment/nginx-config-new.conf"
    echo ""
    echo "3. Update /etc/nginx/sites-available/nutriediet.com with the new configuration"
    echo ""
    echo "4. Test the configuration:"
    echo "   sudo nginx -t"
    echo ""
    echo "5. If test passes, reload nginx:"
    echo "   sudo systemctl reload nginx"
    echo ""
    
    read -p "Press Enter when you have completed the nginx configuration..."
}

run_migrations() {
    print_step "Running database migrations..."
    
    cd "$BACKEND_DIR"
    
    # Check if migrate directory exists
    if [ -d "migrate" ]; then
        # Run migrations programmatically if there's a migrate tool
        if [ -f "migrate/migrate.go" ]; then
            go run migrate/migrate.go
            print_success "Migrations completed"
        else
            print_warning "No migrate tool found. You may need to run migrations manually."
        fi
    else
        print_warning "No migrations directory found"
    fi
}

verify_deployment() {
    print_step "Verifying deployment..."
    
    # Check if Go API is running
    sleep 3
    if curl -s http://localhost:8080/health > /dev/null 2>&1 || curl -s http://localhost:8080 > /dev/null 2>&1; then
        print_success "Go API is responding on port 8080"
    else
        print_warning "Go API may not be responding on port 8080"
    fi
    
    # Check if existing app is still running
    if curl -s http://localhost:2299 > /dev/null; then
        print_success "Existing app still running on port 2299"
    else
        print_warning "Existing app may not be responding on port 2299"
    fi
    
    # Check PM2 status
    echo ""
    pm2 list
    
    echo ""
    print_step "Deployment verification complete!"
    echo ""
    echo "Next steps:"
    echo "1. Access your new app at: https://nutriediet.com/new"
    echo "2. Test the API at: https://nutriediet.com/new/api"
    echo "3. Verify existing site still works at: https://nutriediet.com"
    echo ""
    echo "Useful commands:"
    echo "  - View Go API logs: pm2 logs nutriediet-go-api"
    echo "  - Restart Go API: pm2 restart nutriediet-go-api"
    echo "  - View all PM2 apps: pm2 list"
    echo "  - View nginx logs: sudo tail -f /var/log/nginx/error.log"
}

# =============================================================================
# Main Execution
# =============================================================================

main() {
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║     Nutriediet New App Deployment Script                  ║"
    echo "║     Go Backend + React Frontend to /new subpath           ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    check_user
    check_existing_site
    
    print_warning "This script will:"
    echo "  - Install Go $GO_VERSION (if needed)"
    echo "  - Upgrade Node.js to v$NODE_VERSION (if needed)"
    echo "  - Create new MySQL database: $DB_NAME"
    echo "  - Clone and build your application"
    echo "  - Configure PM2 to run Go backend on port 8080"
    echo "  - Update Nginx configuration"
    echo ""
    
    read -p "Continue with deployment? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
    
    install_go
    upgrade_node
    create_directories
    setup_database
    clone_repos
    build_backend
    configure_backend
    build_frontend
    run_migrations
    setup_pm2
    configure_nginx
    verify_deployment
    
    print_success "Deployment complete! 🎉"
}

# Run main function
main

