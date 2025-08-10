#!/bin/bash

# ViberCode CLI Quick Start Script
# This script helps you get started with ViberCode CLI quickly

set -e

echo "🚀 ViberCode CLI Quick Start"
echo "============================"
echo

# Check if vibercode is installed
if ! command -v vibercode &> /dev/null; then
    echo "❌ ViberCode CLI is not installed."
    echo "💡 Please install it first:"
    echo "   go install github.com/vibercode/cli@latest"
    echo "   # or use the install script: ./install.sh"
    exit 1
fi

echo "✅ ViberCode CLI found!"
echo

# Ask user for project type
echo "🎯 What type of project would you like to create?"
echo "1) Simple Blog API"
echo "2) User Management System"
echo "3) Custom API (interactive)"
echo

read -p "Choose an option (1-3): " choice

case $choice in
    1)
        SCHEMA="examples/schemas/blog-api.json"
        PROJECT_NAME="simple-blog-api"
        echo "📝 Creating Simple Blog API..."
        ;;
    2)
        SCHEMA="examples/schemas/user-management.json"
        PROJECT_NAME="user-management-api"
        echo "👥 Creating User Management System..."
        ;;
    3)
        echo "🎨 Starting interactive API generator..."
        vibercode generate api
        echo "✅ Custom API generated successfully!"
        exit 0
        ;;
    *)
        echo "❌ Invalid option. Please choose 1, 2, or 3."
        exit 1
        ;;
esac

# Check if schema file exists
if [ ! -f "$SCHEMA" ]; then
    echo "❌ Schema file not found: $SCHEMA"
    echo "💡 Please run this script from the vibercode-cli-go directory"
    exit 1
fi

# Ask for project directory
read -p "📁 Enter project directory name [$PROJECT_NAME]: " user_project_name
PROJECT_NAME=${user_project_name:-$PROJECT_NAME}

# Check if directory already exists
if [ -d "$PROJECT_NAME" ]; then
    echo "⚠️  Directory '$PROJECT_NAME' already exists."
    read -p "Do you want to remove it and continue? (y/N): " confirm
    if [[ $confirm =~ ^[Yy]$ ]]; then
        rm -rf "$PROJECT_NAME"
        echo "🗑️  Removed existing directory"
    else
        echo "❌ Aborted."
        exit 1
    fi
fi

# Generate the API
echo "⚡ Generating API..."
vibercode generate api --schema "$SCHEMA" --output "$PROJECT_NAME"

if [ $? -eq 0 ]; then
    echo "✅ API generated successfully!"
    echo
    echo "🎉 Next steps:"
    echo "   cd $PROJECT_NAME"
    echo "   make setup    # Set up dependencies and database"
    echo "   make dev      # Start development server"
    echo
    echo "📖 Check the README.md in your project for detailed instructions."
    echo
    
    # Ask if user wants to start the project
    read -p "🚀 Do you want to start the development server now? (y/N): " start_now
    if [[ $start_now =~ ^[Yy]$ ]]; then
        cd "$PROJECT_NAME"
        echo "📦 Setting up dependencies..."
        make setup
        echo "🌐 Starting development server..."
        make dev
    fi
else
    echo "❌ Failed to generate API. Please check the error messages above."
    exit 1
fi
