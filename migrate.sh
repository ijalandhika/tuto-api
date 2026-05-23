#!/bin/bash

# Usage:
#   ./migrate.sh up       ← apply all migrations
#   ./migrate.sh down     ← rollback one step
#   ./migrate.sh down 3   ← rollback 3 steps
#   ./migrate.sh force 1  ← force set version

 SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.yml"
MIGRATIONS_PATH="$SCRIPT_DIR/migrations"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: $CONFIG_FILE not found. Copy config.example.yml to config.yml first."
    exit 1
fi

DB_URL=$(yq '.database.url' "$CONFIG_FILE" | tr -d '"')

DIRECTION=${1:-up}
STEPS=${2:-}

echo "Running migration: $DIRECTION"
echo "DB: $DB_URL"

if [ -z "$STEPS" ]; then
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" "$DIRECTION"
else
    migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" "$DIRECTION" "$STEPS"
fi