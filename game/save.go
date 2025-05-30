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
	Timestamp  time.Time            `json:"timestamp"`
	Tick       int                  `json:"tick"`
	Age        string               `json:"age"`
	Resources  map[string]float64   `json:"resources"`
	Buildings  map[string]int       `json:"buildings"`
	Villagers  map[string]VillagerInfo `json:"villagers"`
	Stats      *GameStats           `json:"stats"`
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
		Timestamp: time.Now(),
		Tick:      ge.Tick,
		Age:       ge.Age,
		Resources: ge.Resources.GetAll(),
		Buildings: ge.Buildings.GetAll(),
		Villagers: ge.Villagers.GetAll(),
		Stats:     ge.Stats,
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
		return fmt.Errorf("failed to unmarshal save data: %w", err)
	}

	// Restore game state
	ge.Tick = save.Tick
	ge.Age = save.Age

	// Restore resources
	for resource, amount := range save.Resources {
		ge.Resources.resources[resource] = amount
	}

	// Restore buildings
	for building, count := range save.Buildings {
		ge.Buildings.buildings[building] = count
	}

	// Restore villagers
	for vtype, info := range save.Villagers {
		if _, exists := ge.Villagers.villagers[vtype]; exists {
			ge.Villagers.villagers[vtype].Count = info.Count
			for resource, count := range info.Assignment {
				ge.Villagers.villagers[vtype].Assignment[resource] = count
			}
		}
	}
	
	// Restore statistics if available
	if save.Stats != nil {
		ge.Stats = save.Stats
	}

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
