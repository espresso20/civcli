package main

import (
	"fmt"
	"time"

	"github.com/user/civcli/game"
)

func main() {
	fmt.Println("Testing game engine load functionality...")

	// Create a minimal game engine (similar to how main.go does it)
	gameEngine := game.NewGameEngine(nil) // No UI for this test

	// Start the engine
	go func() {
		if err := gameEngine.Start(); err != nil {
			fmt.Printf("Error starting game engine: %v\n", err)
		}
	}()

	// Give it a moment to initialize
	time.Sleep(1 * time.Second)

	fmt.Println("Game engine initialized. Current state:")
	state := gameEngine.GetGameState()
	fmt.Printf("Tick: %d, Age: %s\n", state.Tick, state.Age)
	fmt.Printf("Food: %.1f, Wood: %.1f\n", state.Resources["food"], state.Resources["wood"])

	// Now try to load a save file
	fmt.Println("\nTrying to load 'valid_game' save...")
	err := gameEngine.LoadGame("valid_game")
	if err != nil {
		fmt.Printf("Failed to load game: %v\n", err)
		return
	}

	fmt.Println("Load successful! New state:")
	state = gameEngine.GetGameState()
	fmt.Printf("Tick: %d, Age: %s\n", state.Tick, state.Age)
	fmt.Printf("Food: %.1f, Wood: %.1f\n", state.Resources["food"], state.Resources["wood"])

	fmt.Println("Test completed successfully!")
}
