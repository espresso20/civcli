package game

import (
	"fmt"
	"time"
)

// GameEvent represents a significant event in the game
type GameEvent struct {
	Tick      int       `json:"tick"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"eventType"`
	Message   string    `json:"message"`
}

// GameStats keeps track of game statistics
type GameStats struct {
	Events            []GameEvent       `json:"events"`
	ResourcesGathered map[string]float64 `json:"resourcesGathered"`
	BuildingsBuilt    map[string]int    `json:"buildingsBuilt"`
	VillagersRecruited map[string]int   `json:"villagersRecruited"`
	AgesReached       []string          `json:"agesReached"`
	StartTime         time.Time         `json:"startTime"`
}

// NewGameStats creates a new game stats tracker
func NewGameStats() *GameStats {
	return &GameStats{
		Events:            []GameEvent{},
		ResourcesGathered: make(map[string]float64),
		BuildingsBuilt:    make(map[string]int),
		VillagersRecruited: make(map[string]int),
		AgesReached:       []string{"Stone Age"},
		StartTime:         time.Now(),
	}
}

// AddEvent adds a new event to the game history
func (gs *GameStats) AddEvent(tick int, eventType, message string) {
	event := GameEvent{
		Tick:      tick,
		Timestamp: time.Now(),
		EventType: eventType,
		Message:   message,
	}
	gs.Events = append(gs.Events, event)
}

// AddResourceGathered adds to the total resources gathered
func (gs *GameStats) AddResourceGathered(resource string, amount float64) {
	gs.ResourcesGathered[resource] += amount
}

// AddBuildingBuilt increments the count of buildings built
func (gs *GameStats) AddBuildingBuilt(building string) {
	gs.BuildingsBuilt[building]++
}

// AddVillagerRecruited increments the count of villagers recruited
func (gs *GameStats) AddVillagerRecruited(villagerType string) {
	gs.VillagersRecruited[villagerType]++
}

// AddAgeReached adds a new age reached
func (gs *GameStats) AddAgeReached(age string) {
	// Check if age is already in the list
	for _, a := range gs.AgesReached {
		if a == age {
			return
		}
	}
	gs.AgesReached = append(gs.AgesReached, age)
}

// GetPlayTime returns the time played in hours and minutes
func (gs *GameStats) GetPlayTime() string {
	duration := time.Since(gs.StartTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// GetTotalResourcesGathered returns the total of all resources gathered
func (gs *GameStats) GetTotalResourcesGathered() float64 {
	total := 0.0
	for _, amount := range gs.ResourcesGathered {
		total += amount
	}
	return total
}

// GetTotalBuildingsBuilt returns the total of all buildings built
func (gs *GameStats) GetTotalBuildingsBuilt() int {
	total := 0
	for _, count := range gs.BuildingsBuilt {
		total += count
	}
	return total
}

// GetTotalVillagersRecruited returns the total of all villagers recruited
func (gs *GameStats) GetTotalVillagersRecruited() int {
	total := 0
	for _, count := range gs.VillagersRecruited {
		total += count
	}
	return total
}
