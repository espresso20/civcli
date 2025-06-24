package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GameSave represents a saved game state
type GameSave struct {
	Timestamp      time.Time               `json:"timestamp"`
	Tick           int                     `json:"tick"`
	Age            string                  `json:"age"`
	Resources      map[string]float64      `json:"resources"`
	Buildings      map[string]int          `json:"buildings"`
	Villagers      map[string]VillagerInfo `json:"villagers"`
	Stats          *GameStats              `json:"stats"`
	LastUpdateTime time.Time               `json:"lastUpdateTime"`
}

// SaveGame saves the current game state to a file
func (ge *GameEngine) SaveGame(filename string) error {
	// Create save directory if it doesn't exist
	saveDir := filepath.Join(".", "data", "saves")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// Prepare save data
	save := GameSave{
		Timestamp:      time.Now(),
		Tick:           ge.Tick,
		Age:            ge.Age,
		Resources:      ge.Resources.GetAll(),
		Buildings:      ge.Buildings.GetAll(),
		Villagers:      ge.Villagers.GetAll(),
		Stats:          ge.Stats,
		LastUpdateTime: ge.LastUpdateTime,
	}

	// Marshal to JSON
	saveData, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Write to file
	savePath := filepath.Join(saveDir, filename+".json")
	if err := os.WriteFile(savePath, saveData, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadGame loads a game state from a file
func (ge *GameEngine) LoadGame(filename string) error {
	savePath := filepath.Join(".", "data", "saves", filename+".json")

	// Read the save file
	saveData, err := os.ReadFile(savePath)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal the JSON
	var save GameSave
	if err := json.Unmarshal(saveData, &save); err != nil {
		return fmt.Errorf("failed to unmarshal save data - save file may be corrupted or from an incompatible version: %w", err)
	}

	// Validate essential fields
	if save.Age == "" {
		return fmt.Errorf("save file is missing required 'age' field")
	}
	if save.Resources == nil {
		return fmt.Errorf("save file is missing required 'resources' field")
	}

	// Note: We don't stop the game engine anymore since it can cause the main loop to exit
	// The game state update is atomic enough that we can safely update while running

	// Restore game state safely
	ge.Tick = save.Tick
	ge.Age = save.Age

	// Restore resources (ensure Resources manager exists)
	if ge.Resources == nil {
		ge.Resources = NewResourceManager()
	}
	// Clear and restore resources
	ge.Resources.resources = make(map[string]float64)
	for resource, amount := range save.Resources {
		ge.Resources.resources[resource] = amount
	}

	// Restore buildings (ensure Buildings manager exists)
	if ge.Buildings == nil {
		ge.Buildings = NewBuildingManager()
	}
	// Clear and restore buildings
	ge.Buildings.buildings = make(map[string]int)
	if save.Buildings != nil {
		for building, count := range save.Buildings {
			ge.Buildings.buildings[building] = count
		}
	}

	// Restore villagers (ensure Villagers manager exists)
	if ge.Villagers == nil {
		ge.Villagers = NewVillagerManager()
	}
	// Clear and restore villagers
	ge.Villagers.villagers = make(map[string]*VillagerType)
	if save.Villagers != nil {
		for vtype, info := range save.Villagers {
			// Create new villager entry with safe defaults
			assignment := make(VillagerAssignment)
			if info.Assignment != nil {
				for resource, count := range info.Assignment {
					assignment[resource] = count
				}
			}

			ge.Villagers.villagers[vtype] = &VillagerType{
				Count:      info.Count,
				FoodCost:   0.5, // Default food cost
				Assignment: assignment,
			}
		}
	}

	// Ensure other managers exist
	if ge.Progress == nil {
		ge.Progress = NewProgressManager()
	}
	if ge.Research == nil {
		ge.Research = NewResearchManager()
	}
	// if ge.Library == nil {
	// 	ge.Library = NewLibrarySystem()
	// }

	// Restore or create statistics
	if save.Stats != nil {
		ge.Stats = save.Stats
	} else {
		if ge.Stats == nil {
			ge.Stats = NewGameStats()
		}
	}

	// Restore LastUpdateTime or set to current time
	if !save.LastUpdateTime.IsZero() {
		ge.LastUpdateTime = save.LastUpdateTime
	} else {
		ge.LastUpdateTime = time.Now()
	}

	// No need to restart the engine since we didn't stop it
	return nil
}

// ListSaves returns a list of available save files
func ListSaves() ([]string, error) {
	saveDir := filepath.Join(".", "data", "saves")

	// Check if directory exists
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory
	files, err := os.ReadDir(saveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read save directory: %w", err)
	}

	// Extract filenames
	var saves []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			saves = append(saves, file.Name()[:len(file.Name())-5]) // Remove .json extension
		}
	}

	return saves, nil
}
