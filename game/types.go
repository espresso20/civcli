package game

// GameState represents the essential state of the game that can be accessed by the UI
type GameState struct {
	Age         string
	Tick        int
	Resources   map[string]float64
	Buildings   map[string]int
	Villagers   map[string]VillagerInfo
	VillagerCap int
	Research    struct {
		Current    string
		Progress   float64
		Cost       float64
		Researched []string
	}
	TickDurationSeconds float64  // Add tick duration (seconds per tick) for UI display
	foodSources         []string // List of resources that count as food
	TotalFood           float64  // Total amount of food from all food sources
}

// GetTotalFood returns the sum of all food resources
func (gs *GameState) GetTotalFood() float64 {
	var total float64

	// If we have the food sources list, use it
	if len(gs.foodSources) > 0 {
		for _, source := range gs.foodSources {
			if amount, exists := gs.Resources[source]; exists {
				total += amount
			}
		}
	} else {
		// Default food sources if not explicitly defined
		foodSources := []string{"foraging", "hunting"}
		for _, source := range foodSources {
			if amount, exists := gs.Resources[source]; exists {
				total += amount
			}
		}
	}

	return total
}

// GameStateProvider defines an interface for accessing game state information
// This interface is used by the UI to avoid import cycles
type GameStateProvider interface {
	GetGameState() GameState
}

// Ensure GameEngine implements GameStateProvider
var _ GameStateProvider = (*GameEngine)(nil)

// GetGameState returns the current game state for UI rendering
func (ge *GameEngine) GetGameState() GameState {
	// Get research info
	currentResearch, progress, cost := ge.Research.GetProgress()

	// Get researched technologies
	researchedTechs := ge.Research.GetResearchedTechnologies()
	researched := make([]string, 0, len(researchedTechs))
	for name := range researchedTechs {
		researched = append(researched, name)
	}

	// Create GameState with food sources
	gameState := GameState{
		Age:         ge.Age,
		Tick:        ge.Tick,
		Resources:   ge.Resources.GetAll(),
		Buildings:   ge.Buildings.GetAll(),
		Villagers:   ge.Villagers.GetAll(),
		VillagerCap: ge.Buildings.GetVillagerCapacity(),
		Research: struct {
			Current    string
			Progress   float64
			Cost       float64
			Researched []string
		}{
			Current:    currentResearch,
			Progress:   progress,
			Cost:       cost,
			Researched: researched,
		},
		TickDurationSeconds: ge.TickDuration.Seconds(),   // Pass tick duration to UI
		TotalFood:           ge.Resources.GetTotalFood(), // Pass total food value separately
	}

	return gameState
}
