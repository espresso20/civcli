#!/bin/zsh

# This script builds and runs the CivIdleCli game

# Create bin directory if it doesn't exist
mkdir -p bin

echo "Building CivIdleCli..."
go build -o bin/cividlecli

if [ $? -eq 0 ]; then
    echo "Build successful! Running the game..."
    ./bin/cividlecli
else
    echo "Build failed."
    exit 1
fi
