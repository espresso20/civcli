package game

import (
	"errors"
	"time"
)

// GameEngine represents the main game state and logic
type GameEngine struct {
	Display      DisplayInterface
	Running      bool
	Tick         int
	Age          string
	Resources    *ResourceManager
	Buildings    *BuildingManager
	Villagers    *VillagerManager
	Progress     *ProgressManager
	Research     *ResearchManager
	Library      *LibrarySystem
	Commands     *CommandHandler
	Stats        *GameStats
	TickDuration time.Duration
}

// DisplayInterface defines the interface for the UI display
type DisplayInterface interface {
	ShowIntro()
	ShowHelp(commands map[string]string)
	ShowMessage(message string, style string)
	ShowAgeAdvancement(newAge string)
	DisplayDashboard(state GameState)
	ShowLibraryContent(title, content string)
	ShowLibraryTopicsList(topics map[string]string)
	GetInput() (string, error)
	Stop()
}

// NewGameEngine creates a new game engine
func NewGameEngine(display DisplayInterface) *GameEngine {
	ge := &GameEngine{
		Display:      display,
		Running:      false,
		Tick:         0,
		Age:          "Stone Age",
		TickDuration: 1 * time.Second,
	}

	// Initialize game components
	ge.Resources = NewResourceManager()
	ge.Buildings = NewBuildingManager()
	ge.Villagers = NewVillagerManager()
	ge.Progress = NewProgressManager()
	ge.Research = NewResearchManager()
	ge.Library = NewLibrarySystem()
	ge.Stats = NewGameStats()

	// Initialize game state
	ge.initializeGame()

	// Create command handler (after initializing components)
	ge.Commands = NewCommandHandler(ge)

	return ge
}

// initializeGame sets up the initial game state
func (ge *GameEngine) initializeGame() {
	// Add initial resources
	ge.Resources.Add("food", 20)
	ge.Resources.Add("wood", 15)

	// Start with one villager
	ge.Villagers.Add("villager", 1)
}

// Start starts the main game loop
func (ge *GameEngine) Start() error {
	ge.Running = true
	return ge.mainLoop()
}

// mainLoop is the main game loop
func (ge *GameEngine) mainLoop() error {
	for ge.Running {
		// Display the dashboard
		ge.Display.DisplayDashboard(ge.GetGameState())

		// Get user command
		userInput, err := ge.Display.GetInput()
		if err != nil {
			return errors.New("error getting user input: " + err.Error())
		}

		// Process command
		ge.Commands.Process(userInput)

		// Update game state (process one tick)
		ge.update()
	}
	return nil
}

// update updates the game state for one tick
func (ge *GameEngine) update() {
	ge.Tick++

	// Update resources based on villagers and track statistics
	ge.Villagers.CollectResourcesAndTrack(ge.Resources, ge.Stats)

	// Update buildings
	ge.Buildings.Update(ge.Resources)

	// Update research if there's an active research project
	if techName, completed := ge.Research.ContinueResearch(ge.Resources.Get("knowledge") * 0.1); completed {
		ge.Display.ShowMessage("Research completed: "+techName, "success")
		ge.Stats.AddEvent(ge.Tick, "research_completed", "Completed research on "+techName)

		// Apply any immediate effects from the research
		// For now, this is just handled by the production rate calculation
	}

	// Check for age progression
	newAge := ge.Progress.CheckAdvancement(ge.Resources, ge.Buildings, ge.Age)
	if newAge != "" && newAge != ge.Age {
		ge.Display.ShowAgeAdvancement(newAge)
		ge.Age = newAge

		// Track age advancement in stats
		ge.Stats.AddEvent(ge.Tick, "age_advancement", "Advanced to "+newAge)
		ge.Stats.AddAgeReached(newAge)
	}
}

// Quit quits the game
func (ge *GameEngine) Quit() {
	ge.Display.ShowMessage("Goodbye! Thanks for playing CivIdleCli!", "warning")

	// Sleep briefly to allow the message to be displayed
	time.Sleep(500 * time.Millisecond)

	// Stop the display (which stops the tview application)
	if stopDisplay, ok := ge.Display.(interface{ Stop() }); ok {
		stopDisplay.Stop()
	}

	ge.Running = false
}
