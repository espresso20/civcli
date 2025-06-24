package main

import (
	"fmt"
	"log"
	"os"

	"github.com/user/civcli/game"
)

// DummyDisplay implements DisplayInterface for testing
type DummyDisplay struct{}

func (d *DummyDisplay) DisplayDashboard(state game.GameState) {}
func (d *DummyDisplay) ShowMessage(message, msgType string)   {}
func (d *DummyDisplay) ShowHelp(commands map[string]string)   {}
func (d *DummyDisplay) ShowAgeAdvancement(newAge string)      {}
func (d *DummyDisplay) GetInput() (string, error)             { return "", fmt.Errorf("no input") }
func (d *DummyDisplay) Stop()                                 {}

func main() {
	testLoad()
}

func testLoad() {
	// Test just the game engine loading functionality
	fmt.Println("Testing minimal game load functionality...")

	// Create a new game engine with dummy display
	engine := game.NewGameEngine(&DummyDisplay{})

	// Check what save files exist
	fmt.Println("\nChecking save files...")
	files, err := os.ReadDir("./data/saves")
	if err != nil {
		log.Fatalf("Failed to read saves directory: %v", err)
	}

	fmt.Printf("Found %d files in saves directory:\n", len(files))
	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("  - %s\n", file.Name())
		}
	}

	// Try to load a known good save
	fmt.Println("\nTesting load of early_game save...")
	err = engine.LoadGame("early_game")
	if err != nil {
		log.Fatalf("Failed to load early_game: %v", err)
	}

	fmt.Println("✅ Successfully loaded early_game!")

	// Print some basic info
	state := engine.GetGameState()
	fmt.Printf("Loaded game state:\n")
	fmt.Printf("  Tick: %d\n", state.Tick)
	fmt.Printf("  Age: %s\n", state.Age)
	fmt.Printf("  Resources: %v\n", state.Resources)

	fmt.Println("\n✅ All tests passed! Game loading works correctly in isolation.")
}
