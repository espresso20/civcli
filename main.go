package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/civcli/game"
	"github.com/user/civcli/ui"
)

func main() {
	fmt.Println("CivIdleCli - A Command Line Civilization Builder")

	// Set up signal handling for clean exit
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\nExiting CivIdleCli...")
		os.Exit(0)
	}()

	// Initialize the display
	display := ui.NewDisplay()

	// Initialize the game engine
	gameEngine := game.NewGameEngine(display)

	// Show the introduction
	display.ShowIntro()

	// Start the game
	err := gameEngine.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}
