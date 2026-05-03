#!/bin/bash

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."
    # Assuming Debian-based system as per system check
    sudo apt update && sudo apt install -y golang
else
    echo "Go is already installed: $(go version)"
fi

# Check for .env file
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Prerequisite: Please copy .env-example to .env and fill in the required tokens before running this script."
    exit 1
fi

# Download dependencies
echo "Ensuring dependencies are downloaded..."
go mod download

# Get absolute path of the current directory
SCRIPT_DIR=$(pwd)
GO_PATH=$(which go)

# Check if cron job already exists
(crontab -l 2>/dev/null | grep -q "request_match.go") && { echo "Cron job already exists."; exit 0; }

# Add cron job (9 AM every day)
(crontab -l 2>/dev/null; echo "0 9 * * * cd $SCRIPT_DIR && $GO_PATH run request_match.go >> $SCRIPT_DIR/cron.log 2>&1") | crontab -

echo "Installation complete! Cron job added for 9:00 AM daily."
