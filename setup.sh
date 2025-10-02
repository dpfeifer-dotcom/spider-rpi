#!/bin/bash
set -e

# Frissítés
sudo apt-get update -y

# Git telepítése
echo ">>> Git telepítése..."
sudo apt-get install -y git

# Docker telepítése (docker.io)
echo ">>> Docker telepítése..."
sudo apt-get install -y docker.io docker-compose-plugin

# Docker engedélyezése induláskor
sudo systemctl enable docker
sudo systemctl start docker

# Hozzáadjuk az aktuális usert a docker csoporthoz
sudo usermod -aG docker $USER

# Docker Compose fájl letöltése
COMPOSE_URL="https://raw.githubusercontent.com/dpfeifer-dotcom/spider-rpi/main/docker-compose.yaml"
TARGET_DIR="$HOME/spider-rpi"

echo ">>> Docker Compose fájl letöltése: $COMPOSE_URL"
mkdir -p "$TARGET_DIR"
cd "$TARGET_DIR"
curl -L "$COMPOSE_URL" -o docker-compose.yaml

# Docker Compose indítása
echo ">>> Docker Compose stack indítása..."
docker compose up -d

echo ">>> Kész! A stack fut."
echo ">>> FONTOS: Jelentkezz ki és be újra, hogy a docker csoportos jogosultság érvényesüljön."
