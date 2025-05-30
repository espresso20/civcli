# CivIdleCli - Command Line Civilization Builder Game

CivIdleCli is a text-based clicker-style progression game where you build and manage your own civilization through different ages, all within your terminal.

## Overview

In CivIdleCli, you start with a small tribe in the Stone Age. Your goal is to gather resources, build structures, recruit villagers, and advance through multiple ages of civilization.

## Features

- **Resource Management**: Gather and manage resources like food, wood, stone, gold, and knowledge
- **Villager System**: Recruit villagers and assign them to different tasks
- **Building System**: Construct various buildings that provide bonuses and unlock new capabilities
- **Age Progression**: Advance through different ages, from Stone Age to Modern Age
- **Command-based Interface**: Simple text commands with auto-completion
- **Rich Terminal UI**: Colorful and informative terminal interface
- **Distributable Binary**: Built with Go, can be distributed as a single binary file

## Installation

### Pre-built Binary

Download the latest release for your platform from the Releases page.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/cividlecli.git
cd cividlecli

# Build the project
go build -o cividlecli

# Run the game
./cividlecli
```

## How to Play

### Basic Commands

- `help` - Display available commands
- `gather <resource> <count>` - Assign villagers to gather resources
- `build <building>` - Build a structure
- `recruit <villager_type> <count>` - Recruit new villagers
- `assign <villager_type> <resource> <count>` - Assign villagers to tasks
- `status` - Show detailed status of your civilization
- `buildings` - List available buildings and their costs
- `quit` - Exit the game

## Game Progression

1. Start in the Stone Age with basic resources
2. Build huts to increase population capacity
3. Assign villagers to gather resources
4. Build farms for food production
5. Advance to the Bronze Age
6. Continue expanding and advancing through ages
7. Unlock new buildings, resources, and villager types
8. Reach the Modern Age and build an advanced civilization

## Requirements

- To build from source: Go 1.21 or higher
