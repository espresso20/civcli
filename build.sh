#!/bin/bash

# Build script for CivIdleCli

echo "Building CivIdleCli..."

# Build for current platform
go build -o cividlecli

if [ $? -eq 0 ]; then
    echo "Build successful! Run ./cividlecli to start the game."
else
    echo "Build failed."
    exit 1
fi
