#!/bin/bash

# GOTS Template Setup Script
# This script helps you set up a new project from the GOTS Template

set -e

echo "🚀 GOTS Template Setup"
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
        echo "Error: $field cannot be empty"
        exit 1
    fi
}

# Function to validate Go module name format
validate_module_name() {
    local module="$1"
    
    if [[ ! "$module" =~ ^[a-zA-Z0-9._-]+(/[a-zA-Z0-9._-]+)*$ ]]; then
        echo "Error: Invalid Go module name format"
        echo "Module name should be in format: domain.com/username/project-name"
        exit 1
    fi
}

# Get current directory name as default project name
current_dir=$(basename "$(pwd)")

echo "📝 Project Configuration"
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
echo "📋 Configuration Summary"
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
echo "🔧 Setting up your project..."

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

# Step 7: Create .env file from .env.sample
echo "🔧 Creating .env file from .env.sample..."
if [ -f .env.sample ]; then
    if [ ! -f .env ]; then
        cp .env.sample .env
        echo "✅ Created .env file. Please update it with your configuration values."
    else
        echo "⚠️  .env file already exists, skipping creation."
    fi
else
    echo "⚠️  .env.sample not found, skipping .env creation."
fi

# Step 8: Update README.md
echo "📝 Updating README.md..."
if [ -f README.md ]; then
    sed -i.bak "s|# GOTS Template|# $PROJECT_NAME|g" README.md
    sed -i.bak "s|A production-ready full-stack application template|$DESCRIPTION|g" README.md
    sed -i.bak "s|Run \`./setup-template.sh\` to customize for your project|Project has been set up and customized|g" README.md
    sed -i.bak "s|make run|./$BINARY_NAME run|g" README.md
    rm README.md.bak
    echo "✅ Updated README.md with project information."
fi

# Step 9: Update main.go CLI descriptions
echo "🔧 Updating CLI descriptions in main.go..."
if [ -f main.go ]; then
    sed -i.bak "s|GOTS Template CLI|$PROJECT_NAME CLI|g" main.go
    sed -i.bak "s|A CLI for managing and running the GOTS Template application|A CLI for managing and running $PROJECT_NAME|g" main.go
    sed -i.bak "s|GOTS Template|$PROJECT_NAME|g" main.go
    sed -i.bak "s|Run the GOTS Template application|Run the $PROJECT_NAME application|g" main.go
    sed -i.bak "s|Starts the GOTS Template application|Starts the $PROJECT_NAME application|g" main.go
    rm main.go.bak
    echo "✅ Updated CLI descriptions in main.go."
fi

# Step 10: Update database filename
echo "🗄️  Updating database configuration..."
if [ -f db/db.go ]; then
    sed -i.bak "s|gots-template.db|$PROJECT_NAME.db|g" db/db.go
    rm db/db.go.bak
    echo "✅ Updated database filename to $PROJECT_NAME.db."
fi

# Step 11: Clean git history (optional)
echo
read -p "Do you want to reinitialize the git repository? This will remove all commit history. (y/N): " git_reinit
if [[ "$git_reinit" =~ ^[Yy]$ ]]; then
    echo "🔄 Reinitializing git repository..."
    rm -rf .git
    git init
    git add .
    git commit -m "Initial commit: $PROJECT_NAME based on GOTS Template"
    echo "✅ Git repository reinitialized"
else
    echo "📝 Git history preserved. You may want to update the remote origin manually."
fi

# Step 12: Run go mod tidy
echo "📦 Updating Go dependencies..."
go mod tidy

echo
echo "🎉 Project setup complete!"
echo
echo "Next Steps:"
echo "1. Configure your OAuth providers (GitHub, Google) if needed"
echo "2. Update the .env file with your configuration"
echo "3. Build and run your application:"
echo "   make build"
echo "   ./$BINARY_NAME run"
echo
echo "5. Access the web application at: http://localhost:8080"
echo
echo "Happy coding! 🚀"
