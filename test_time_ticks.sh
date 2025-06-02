#!/bin/bash

# Test script for time-based tick system in CivIdleCli
# This script demonstrates how the game processes ticks based on elapsed time

echo "CivIdleCli Time-Based Tick System Test"
echo "====================================="
echo ""
echo "This test will demonstrate how the game processes ticks based on elapsed real-world time."
echo "The test will simulate opening the game, waiting, and seeing how many ticks have passed."
echo ""
echo "Current TickDuration setting: 5 seconds per tick"
echo ""

# Compile the game to ensure we have the latest version
echo "Compiling the game..."
go build -o civcli

# First run - should process just 1 tick
echo ""
echo "Test 1: Initial run (should process 1 tick)"
echo "-----------------------------------------"
./civcli test_ticks 1

# Short wait - should process a few ticks
echo ""
echo "Test 2: Short wait (5 seconds, should process ~1 tick)"
echo "---------------------------------------------------"
sleep 5
./civcli test_ticks 2

# Medium wait - should process more ticks
echo ""
echo "Test 3: Medium wait (15 seconds, should process ~3 ticks)"
echo "------------------------------------------------------"
sleep 15
./civcli test_ticks 3

# Longer wait - should process even more ticks
echo ""
echo "Test 4: Longer wait (30 seconds, should process ~6 ticks)"
echo "-------------------------------------------------------"
sleep 30
./civcli test_ticks 4

echo ""
echo "Test completed! If the game processed the expected number of ticks for each test,"
echo "then the time-based tick system is working correctly."
echo ""
echo "Note: The exact number of ticks may vary slightly due to system timing variations."
