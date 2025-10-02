#!/bin/bash
set -e

# Frissítsük a csomaglistát
sudo apt-get update -y

# Git telepítése
echo ">>> Git telepítése..."
sudo apt-get install -y git

# Docker telepítése
echo ">>> Docker telepítése..."
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release


sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Docker engedélyezése induláskor
sudo systemctl enable docker
sudo systemctl start docker

# Hozzáadjuk az aktuális usert a docker csoporthoz (hogy ne kelljen mindig sudo)
sudo usermod -aG docker $USER
# Letöltjük a docker-compose.yaml fájlt
COMPOSE_URL="https://raw.githubusercontent.com/dpfeifer-dotcom/spider-rpi/docker-compose.yaml"
TARGET_DIR="$HOME/spider-rpi"

echo ">>> Docker Compose fájl letöltése: $COMPOSE_URL"
mkdir -p "$TARGET_DIR"
cd "$TARGET_DIR"
curl -L "$COMPOSE_URL" -o docker-compose.yaml

# Docker Compose indítása
echo ">>> Docker Compose stack indítása..."
sudo docker compose up -d

echo ">>> Kész! A stack fut."
echo ">>> FONTOS: Jelentkezz ki és be újra, hogy a docker csoportos jogosultság érvényesüljön."
