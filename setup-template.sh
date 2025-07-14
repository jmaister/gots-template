#!/bin/bash

# GOTS Template Setup Script
# This script helps you set up a new project from the GOTS Template

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 GOTS Template Setup${NC}"
echo "This script will help you set up a new project from the GOTS Template."
echo

# Function to prompt for input with default value
prompt_with_default() {
    local prompt="$1"
    local default="$2"
    local result
    
    if [ -n "$default" ]; then
        read -p "$prompt [$default]: " result
        echo "${result:-$default}"
    else
        read -p "$prompt: " result
        echo "$result"
    fi
}

# Function to validate input
validate_not_empty() {
    local value="$1"
    local field="$2"
    
    if [ -z "$value" ]; then
        echo -e "${RED}Error: $field cannot be empty${NC}"
        exit 1
    fi
}

# Function to validate Go module name format
validate_module_name() {
    local module="$1"
    
    if [[ ! "$module" =~ ^[a-zA-Z0-9._-]+(/[a-zA-Z0-9._-]+)*$ ]]; then
        echo -e "${RED}Error: Invalid Go module name format${NC}"
        echo "Module name should be in format: domain.com/username/project-name"
        exit 1
    fi
}

# Get current directory name as default project name
current_dir=$(basename "$(pwd)")

echo -e "${YELLOW}📝 Project Configuration${NC}"
echo

# Collect project information
PROJECT_NAME=$(prompt_with_default "Project name" "$current_dir")
validate_not_empty "$PROJECT_NAME" "Project name"

BINARY_NAME=$(prompt_with_default "Binary executable name" "${PROJECT_NAME}")
validate_not_empty "$BINARY_NAME" "Binary name"

MODULE_NAME=$(prompt_with_default "Go module name (e.g., github.com/username/project-name)" "github.com/username/$PROJECT_NAME")
validate_not_empty "$MODULE_NAME" "Module name"
validate_module_name "$MODULE_NAME"

DESCRIPTION=$(prompt_with_default "Project description" "A full-stack application built with GOTS")

AUTHOR_NAME=$(prompt_with_default "Author name" "$(git config user.name 2>/dev/null || echo '')")
AUTHOR_EMAIL=$(prompt_with_default "Author email" "$(git config user.email 2>/dev/null || echo '')")

echo
echo -e "${YELLOW}📋 Configuration Summary${NC}"
echo "Project Name: $PROJECT_NAME"
echo "Binary Name: $BINARY_NAME"
echo "Module Name: $MODULE_NAME"
echo "Description: $DESCRIPTION"
echo "Author: $AUTHOR_NAME <$AUTHOR_EMAIL>"
echo

read -p "Continue with this configuration? (y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 0
fi

echo
echo -e "${GREEN}🔧 Setting up your project...${NC}"

# Step 1: Update go.mod
echo "📦 Updating Go module..."
sed -i.bak "s|module github.com/jmaister/gots-template|module $MODULE_NAME|g" go.mod
rm go.mod.bak

# Step 2: Update all Go import statements
echo "🔄 Updating import statements..."
find . -name "*.go" -type f -exec sed -i.bak "s|github.com/jmaister/gots-template|$MODULE_NAME|g" {} \;
find . -name "*.go.bak" -delete

# Step 3: Update Makefile
echo "⚙️  Updating Makefile..."
sed -i.bak "s|PROJECT_NAME := gots-template|PROJECT_NAME := $PROJECT_NAME|g" Makefile
sed -i.bak "s|BINARY_NAME := gots|BINARY_NAME := $BINARY_NAME|g" Makefile
rm Makefile.bak

# Step 4: Clean up development artifacts
echo "🧹 Cleaning up development files..."
rm -rf dist/
rm -rf webapp/node_modules/
rm -rf webapp/dist/

# Step 5: Update package.json in webapp
echo "📱 Updating webapp configuration..."
if [ -f webapp/package.json ]; then
    sed -i.bak "s|\"name\": \"webapp\"|\"name\": \"$PROJECT_NAME-webapp\"|g" webapp/package.json
    rm webapp/package.json.bak
fi

# Step 6: Update OpenAPI specification
echo "📋 Updating OpenAPI specification..."
sed -i.bak "s|title: GOTS API|title: $PROJECT_NAME API|g" api/openapi-spec.yaml
sed -i.bak "s|description: API specification for the GOTS Template application|description: API specification for $PROJECT_NAME|g" api/openapi-spec.yaml
rm api/openapi-spec.yaml.bak

# Step 7: Clean git history (optional)
echo
read -p "Do you want to reinitialize the git repository? This will remove all commit history. (y/N): " git_reinit
if [[ "$git_reinit" =~ ^[Yy]$ ]]; then
    echo "🔄 Reinitializing git repository..."
    rm -rf .git
    git init
    git add .
    git commit -m "Initial commit: $PROJECT_NAME based on GOTS Template"
    echo -e "${GREEN}✅ Git repository reinitialized${NC}"
else
    echo "📝 Git history preserved. You may want to update the remote origin manually."
fi

# Step 8: Run go mod tidy
echo "📦 Updating Go dependencies..."
go mod tidy

echo
echo -e "${GREEN}🎉 Project setup complete!${NC}"
echo
echo -e "${BLUE}Next Steps:${NC}"
echo "1. Update the README.md file with your project-specific information"
echo "2. Configure your OAuth providers (GitHub, Google) if needed"
echo "3. Update the .env file with your configuration"
echo "4. Build and run your application:"
echo "   ${YELLOW}make build${NC}"
echo "   ${YELLOW}make run${NC}"
echo
echo "5. Access the web application at: http://localhost:8080"
echo
echo -e "${GREEN}Happy coding! 🚀${NC}"
