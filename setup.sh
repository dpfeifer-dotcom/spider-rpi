#!/bin/bash
clear
set -e

# Colors
GREEN="\033[1;32m"
BLUE="\033[1;34m"
NC="\033[0m" # No color

function info() {
    echo -e "${BLUE}>>> $1${NC}"
}

function success() {
    echo -e "${GREEN}✔ $1${NC}"
}

# Braille spinner frames
frames=(⠋ ⠙ ⠹ ⠸ ⠼ ⠴ ⠦ ⠧ ⠇ ⠏)
frame_count=${#frames[@]}

# Function to run a command with braille spinner
function run_with_spinner() {
    local msg="$1"
    local cmd="$2"
    info "$msg"

    # Run the command in background
    bash -c "$cmd" > /dev/null 2>&1 &
    pid=$!

    local i=0
    while kill -0 $pid 2>/dev/null; do
        printf "\r${frames[$((i % frame_count))]} $msg..."
        i=$((i+1))
        sleep 0.1
    done

    wait $pid
    # Clear spinner line
    printf "\r\033[K"
    success "$msg done."
}

# Update packages
run_with_spinner "Updating package list" "sudo apt-get update -y"

# Install Git
run_with_spinner "Installing Git" "sudo apt-get install -y git"

# Install Docker
run_with_spinner "Installing Docker (docker.io + docker-compose)" "sudo apt-get install -y docker.io docker-compose"

# Enable and start Docker service
run_with_spinner "Enabling and starting Docker service" "sudo systemctl enable docker && sudo systemctl start docker"

# Download Docker Compose file
COMPOSE_URL="https://raw.githubusercontent.com/dpfeifer-dotcom/spider-rpi/main/docker-compose.yaml"
TARGET_DIR="$HOME/spider-rpi"
run_with_spinner "Downloading Docker Compose file" "mkdir -p \"$TARGET_DIR\" && cd \"$TARGET_DIR\" && curl -s -L \"$COMPOSE_URL\" -o docker-compose.yaml"

# Start Docker Compose stack – full output visible
info "Starting Docker Compose stack..."
sudo docker compose -f "$HOME/spider-rpi/docker-compose.yaml" up -d
success "Docker Compose stack is running."

info "Done! Please log out and log back in for Docker group permissions to take effect."

