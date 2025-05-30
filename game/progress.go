package game

// ProgressManager handles age progression and advancement
type ProgressManager struct {
	ages            []string
	ageRequirements map[string]AgeRequirement
	ageUnlocks      map[string]AgeUnlock
}

// AgeRequirement defines what's needed to advance to an age
type AgeRequirement struct {
	Resources map[string]float64
	Buildings map[string]int
}

// AgeUnlock defines what gets unlocked in an age
type AgeUnlock struct {
	Buildings []string
	Resources []string
	Villagers []string
}

// NewProgressManager creates a new progress manager
func NewProgressManager() *ProgressManager {
	pm := &ProgressManager{
		ages: []string{
			"Stone Age",
			"Bronze Age",
			"Iron Age",
			"Medieval Age",
			"Renaissance Age",
			"Industrial Age",
			"Modern Age",
		},
		ageRequirements: map[string]AgeRequirement{
			"Bronze Age": {
				Resources: map[string]float64{"stone": 50, "food": 100},
				Buildings: map[string]int{"hut": 3, "farm": 2},
			},
			"Iron Age": {
				Resources: map[string]float64{"stone": 100, "wood": 150, "knowledge": 20},
				Buildings: map[string]int{"mine": 2, "lumber_mill": 2},
			},
			"Medieval Age": {
				Resources: map[string]float64{"stone": 200, "wood": 250, "gold": 50, "knowledge": 50},
				Buildings: map[string]int{"market": 1, "library": 1},
			},
			"Renaissance Age": {
				Resources: map[string]float64{"gold": 150, "knowledge": 100},
				Buildings: map[string]int{"library": 3, "market": 2},
			},
			"Industrial Age": {
				Resources: map[string]float64{"gold": 300, "knowledge": 200},
				Buildings: map[string]int{"library": 5, "market": 4},
			},
			"Modern Age": {
				Resources: map[string]float64{"gold": 500, "knowledge": 400},
				Buildings: map[string]int{"library": 8, "market": 6},
			},
		},
		ageUnlocks: map[string]AgeUnlock{
			"Stone Age": {
				Buildings: []string{"hut", "farm"},
				Resources: []string{"food", "wood"},
			},
			"Bronze Age": {
				Buildings: []string{"lumber_mill", "mine"},
				Resources: []string{"stone"},
			},
			"Iron Age": {
				Buildings: []string{"market", "library"},
				Resources: []string{"gold", "knowledge"},
			},
			"Medieval Age": {
				Villagers: []string{"scholar"},
			},
			"Renaissance Age": {},
			"Industrial Age":  {},
			"Modern Age":      {},
		},
	}
	return pm
}

// GetCurrentAgeIndex returns the index of the current age
func (pm *ProgressManager) GetCurrentAgeIndex(currentAge string) int {
	for i, age := range pm.ages {
		if age == currentAge {
			return i
		}
	}
	return 0
}

// GetNextAge returns the next age after the current one
func (pm *ProgressManager) GetNextAge(currentAge string) string {
	currentIndex := pm.GetCurrentAgeIndex(currentAge)
	if currentIndex < len(pm.ages)-1 {
		return pm.ages[currentIndex+1]
	}
	return ""
}

// CheckAdvancement checks if player can advance to the next age
func (pm *ProgressManager) CheckAdvancement(resources *ResourceManager, buildings *BuildingManager, currentAge string) string {
	nextAge := pm.GetNextAge(currentAge)
	if nextAge == "" {
		return currentAge // Already at the final age
	}

	// Get requirements for the next age
	requirements, exists := pm.ageRequirements[nextAge]
	if !exists {
		return currentAge
	}

	// Check resource requirements
	for resource, amount := range requirements.Resources {
		if resources.Get(resource) < amount {
			return currentAge
		}
	}

	// Check building requirements
	for building, count := range requirements.Buildings {
		if buildings.GetCount(building) < count {
			return currentAge
		}
	}

	// All requirements met, advance to next age
	return nextAge
}

// GetUnlocks returns content unlocked at a specific age
func (pm *ProgressManager) GetUnlocks(age string) AgeUnlock {
	if unlock, exists := pm.ageUnlocks[age]; exists {
		return unlock
	}
	return AgeUnlock{}
}

// GetRequirements returns the requirements for a specific age
func (pm *ProgressManager) GetRequirements(age string) AgeRequirement {
	if req, exists := pm.ageRequirements[age]; exists {
		return req
	}
	return AgeRequirement{}
}

// GetAllAges returns all the ages
func (pm *ProgressManager) GetAllAges() []string {
	return pm.ages
}
