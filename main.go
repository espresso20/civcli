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

	// Show the splash screen and exit if the user quits
	if err := ui.ShowSplashScreen(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running splash screen: %v\n", err)
		os.Exit(1)
	}

	// Now create the main display for the game
	display := ui.NewDisplay()

	// Initialize the game engine with the display
	gameEngine := game.NewGameEngine(display)

	// Create a command handler and connect it to the game engine
	commandHandler := game.NewCommandHandler(gameEngine)
	gameEngine.Commands = commandHandler

	// Run the game engine in a goroutine
	go func() {
		// Start the game engine
		if err := gameEngine.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
			os.Exit(1)
		}
	}()

	// Handle OS signals
	go func() {
		<-sigs
		fmt.Println("\nExiting CivIdleCli...")
		display.Stop()
		os.Exit(0)
	}()

	// Run the UI program (this blocks until the program exits)
	if err := display.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
		os.Exit(1)
	}
}
