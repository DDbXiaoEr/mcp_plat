#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BIN_DIR="$SCRIPT_DIR/bin"

cd "$SCRIPT_DIR"

for prog in "$BIN_DIR"/*; do
    if [ -x "$prog" ] && [ ! -d "$prog" ]; then
        echo "Starting $(basename "$prog")..."
        "$prog" &
    fi
done

echo "All services started."
wait
