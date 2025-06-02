package game

// BuildingManager handles building construction and effects
type BuildingManager struct {
	buildings           map[string]int
	buildingCosts       map[string]map[string]float64
	buildingEffects     map[string]map[string]float64
	buildingRateBonuses map[string]map[string]map[string]float64 // building -> villagerType -> resource -> bonus percentage
}

// NewBuildingManager creates a new building manager
func NewBuildingManager() *BuildingManager {
	bm := &BuildingManager{
		buildings: map[string]int{
			"hut":         0,
			"farm":        0,
			"lumber_mill": 0,
			"mine":        0,
			"market":      0,
			"library":     0,
		},
		buildingCosts: map[string]map[string]float64{
			"hut":         {"wood": 20},
			"farm":        {"wood": 100, "stone": 50, "food": 100},
			"lumber_mill": {"wood": 100, "stone": 300},
			"mine":        {"wood": 100, "stone": 400},
			"market":      {"wood": 200, "stone": 200, "gold": 100},
			"library":     {"wood": 400, "stone": 200, "knowledge": 100},
		},
		buildingEffects: map[string]map[string]float64{
			"hut":         {"villager_capacity": 2},
			"farm":        {"food": 3.5},  // Increased from 2 to improve food production
			"lumber_mill": {"wood": 2},
			"mine":        {"stone": 1, "gold": 0.2},
			"market":      {"gold": 0.5},
			"library":     {"knowledge": 0.5},
		},
		buildingRateBonuses: map[string]map[string]map[string]float64{
			"farm": {
				"villager": {"food": 0.08}, // Increased from 0.05 to improve food gathering efficiency
			},
			"lumber_mill": {
				"villager": {"wood": 0.1}, // Villagers get +10% wood gathering rate per lumber mill
			},
			"mine": {
				"villager": {"stone": 0.05, "gold": 0.05}, // Villagers get +5% stone and gold gathering rate per mine
			},
			"market": {
				"villager": {"gold": 0.1}, // Villagers get +10% gold gathering rate per market
			},
			"library": {
				"scholar":  {"knowledge": 0.15}, // Scholars get +15% knowledge gathering rate per library
				"villager": {"knowledge": 0.02}, // Villagers get +2% knowledge gathering rate per library
			},
		},
	}
	return bm
}

// Add adds new buildings
func (bm *BuildingManager) Add(building string, count int) bool {
	if _, exists := bm.buildings[building]; exists {
		bm.buildings[building] += count
		return true
	}
	return false
}

// Remove removes buildings
func (bm *BuildingManager) Remove(building string, count int) bool {
	if _, exists := bm.buildings[building]; exists && bm.buildings[building] >= count {
		bm.buildings[building] -= count
		return true
	}
	return false
}

// GetCount returns the count of a specific building
func (bm *BuildingManager) GetCount(building string) int {
	return bm.buildings[building]
}

// GetAll returns all buildings and their counts
func (bm *BuildingManager) GetAll() map[string]int {
	return bm.buildings
}

// GetCost returns the cost to build a specific building
func (bm *BuildingManager) GetCost(building string) map[string]float64 {
	if cost, exists := bm.buildingCosts[building]; exists {
		return cost
	}
	return nil
}

// GetEffect returns the effect of a specific building
func (bm *BuildingManager) GetEffect(building string) map[string]float64 {
	if effect, exists := bm.buildingEffects[building]; exists {
		return effect
	}
	return nil
}

// CanBuild checks if we can build a specific building
func (bm *BuildingManager) CanBuild(building string, resources *ResourceManager) bool {
	costs, exists := bm.buildingCosts[building]
	if !exists {
		return false
	}

	// Check if we have enough resources
	for resource, amount := range costs {
		if !resources.Has(resource, amount) {
			return false
		}
	}

	return true
}

// Build builds a new building
func (bm *BuildingManager) Build(building string, resources *ResourceManager) bool {
	if !bm.CanBuild(building, resources) {
		return false
	}

	// Spend resources
	for resource, amount := range bm.buildingCosts[building] {
		resources.Remove(resource, amount)
	}

	// Add the building
	bm.Add(building, 1)
	return true
}

// Update updates resources based on building effects
func (bm *BuildingManager) Update(resources *ResourceManager) {
	for building, count := range bm.buildings {
		if count > 0 {
			if effects, exists := bm.buildingEffects[building]; exists {
				for resource, amount := range effects {
					if resource != "villager_capacity" {
						// Only add direct resource production here, not collection rate bonuses
						resources.Add(resource, amount*float64(count))
					}
				}
			}
		}
	}
}

// GetVillagerCapacity calculates total villager capacity from buildings
func (bm *BuildingManager) GetVillagerCapacity() int {
	capacity := 1 // Start with capacity for 1 villager

	for building, count := range bm.buildings {
		if count > 0 {
			if effects, exists := bm.buildingEffects[building]; exists {
				if cap, hasCapacity := effects["villager_capacity"]; hasCapacity {
					capacity += int(cap * float64(count))
				}
			}
		}
	}

	return capacity
}

// GetCollectionRateBonus returns the collection rate bonus for a specific villager type and resource
func (bm *BuildingManager) GetCollectionRateBonus(villagerType, resource string) float64 {
	totalBonus := 0.0

	// Check each building for applicable bonuses
	for building, count := range bm.buildings {
		if count > 0 {
			// Check if this building provides bonuses for this villager type
			if villagerBonuses, exists := bm.buildingRateBonuses[building]; exists {
				if resourceBonuses, hasVillagerType := villagerBonuses[villagerType]; hasVillagerType {
					if bonus, hasResource := resourceBonuses[resource]; hasResource {
						// Apply bonus for each building of this type
						totalBonus += bonus * float64(count)
					}
				}
			}
		}
	}

	return totalBonus
}
