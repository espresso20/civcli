package game

import (
	"fmt"
	"time"
)

// GameEngine represents the main game state and logic
type GameEngine struct {
	Display   DisplayInterface
	Running   bool
	Tick      int
	Age       string
	Resources *ResourceManager
	Buildings *BuildingManager
	Villagers *VillagerManager
	Progress  *ProgressManager
	Research  *ResearchManager
	// Library        *LibrarySystem
	Commands       *CommandHandler
	Stats          *GameStats
	TickDuration   time.Duration
	LastUpdateTime time.Time
	RefreshRate    time.Duration // How often to refresh the UI
	stopRefresh    chan bool     // Channel to signal stopping the UI refresh
}

// DisplayInterface defines the interface for the UI display
type DisplayInterface interface {
	ShowHelp(commands map[string]string)
	ShowMessage(message string, style string)
	ShowAgeAdvancement(newAge string)
	DisplayDashboard(state GameState)
	GetInput() (string, error)
	Stop()
}

// NewGameEngine creates a new game engine
func NewGameEngine(display DisplayInterface) *GameEngine {
	ge := &GameEngine{
		Display:        display,
		Running:        false,
		Tick:           0,
		Age:            "Stone Age",
		TickDuration:   5 * time.Second,
		LastUpdateTime: time.Now(),
		RefreshRate:    5 * time.Second, // Match refresh rate to tick duration
		stopRefresh:    make(chan bool), // Initialize the stop channel
	}

	// Initialize game components
	ge.Resources = NewResourceManager()
	ge.Buildings = NewBuildingManager()
	ge.Villagers = NewVillagerManager()
	ge.Progress = NewProgressManager()
	ge.Research = NewResearchManager()
	// ge.Library = NewLibrarySystem()
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
	ge.Resources.Add("wood", 20)

	// Start with one villager
	ge.Villagers.Add("villager", 1)
}

// Start initializes and starts the game engine
func (ge *GameEngine) Start() error {
	ge.Running = true

	// Initialize all subsystems if they haven't been already
	if ge.Resources == nil {
		ge.Resources = NewResourceManager()
	}
	if ge.Buildings == nil {
		ge.Buildings = NewBuildingManager()
	}
	if ge.Villagers == nil {
		ge.Villagers = NewVillagerManager()
	}
	if ge.Progress == nil {
		ge.Progress = NewProgressManager()
	}
	if ge.Research == nil {
		ge.Research = NewResearchManager()
	}
	// if ge.Library == nil {
	// 	ge.Library = NewLibrarySystem()
	// }
	if ge.Stats == nil {
		ge.Stats = NewGameStats()
	}
	if ge.Commands == nil {
		ge.Commands = NewCommandHandler(ge)
	}

	ge.RefreshRate = 500 * time.Millisecond
	ge.stopRefresh = make(chan bool)

	// Start the UI refresh goroutine
	go ge.refreshUILoop()

	// Run the main game loop
	err := ge.mainLoop()

	// Signal UI to stop refreshing
	select {
	case <-ge.stopRefresh:
		// Channel already closed, do nothing
	default:
		close(ge.stopRefresh)
	}

	// Return any error from the main loop
	return err
}

// refreshUILoop updates the UI at regular intervals
func (ge *GameEngine) refreshUILoop() {
	ticker := time.NewTicker(ge.RefreshRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Update the game state and refresh the UI
			gameState := ge.GetGameState()
			ge.Display.DisplayDashboard(gameState)
		case <-ge.stopRefresh:
			return
		}
	}
}

// mainLoop is the main game loop
func (ge *GameEngine) mainLoop() error {
	for ge.Running {
		// Get user command (non-blocking)
		userInput, err := ge.Display.GetInput()
		if err != nil {
			// No input ready, continue loop (allows UI updates to happen)
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Process command (ticks are now handled by the refresh loop)
		ge.Commands.Process(userInput)
	}
	return nil
}

// calculateElapsedTicks calculates how many ticks have passed since the last update
func (ge *GameEngine) calculateElapsedTicks() int {
	now := time.Now()
	elapsedTime := now.Sub(ge.LastUpdateTime)
	elapsedTicks := int(elapsedTime.Seconds() / ge.TickDuration.Seconds())

	// Ensure we process at least one tick when the player interacts
	if elapsedTicks < 1 {
		elapsedTicks = 1
	}

	// Update the last update time
	ge.LastUpdateTime = now

	return elapsedTicks
}

// updateSingleTick processes a single tick of game time
func (ge *GameEngine) updateSingleTick() {
	ge.Tick++

	// Update resources based on villagers and track statistics
	ge.Villagers.CollectResourcesAndTrack(ge.Resources, ge.Stats, ge.Buildings)

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

// updateMultipleTicks processes multiple ticks at once
func (ge *GameEngine) updateMultipleTicks(tickCount int) {
	// For large tick counts (e.g. after long absence), limit to a reasonable number
	// to prevent excessive processing and resource accumulation
	maxTicksAtOnce := 1000 // Process at most 1000 ticks at once
	if tickCount > maxTicksAtOnce {
		ge.Display.ShowMessage(fmt.Sprintf("Processing %d ticks (capped from %d)...", maxTicksAtOnce, tickCount), "info")
		tickCount = maxTicksAtOnce
	} else if tickCount > 10 {
		ge.Display.ShowMessage(fmt.Sprintf("Processing %d ticks...", tickCount), "info")
	}

	for i := 0; i < tickCount; i++ {
		ge.updateSingleTick()
	}
}

// Quit quits the game
func (ge *GameEngine) Quit() {
	ge.Display.ShowMessage("Goodbye! Thanks for playing CivIdleCli!", "warning")

	// Signal the refresh loop to stop
	close(ge.stopRefresh)

	// Sleep briefly to allow the message to be displayed
	time.Sleep(500 * time.Millisecond)

	// Stop the display (which stops the tview application)
	if stopDisplay, ok := ge.Display.(interface{ Stop() }); ok {
		stopDisplay.Stop()
	}

	ge.Running = false
}
