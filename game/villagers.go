package game

// VillagerAssignment represents assignment of villagers to tasks
type VillagerAssignment map[string]int

// VillagerType represents a type of villager with its properties
type VillagerType struct {
	Count      int
	FoodCost   float64
	Assignment VillagerAssignment
}

// VillagerManager handles villager creation and assignment
type VillagerManager struct {
	villagers map[string]*VillagerType
}

// NewVillagerManager creates a new villager manager
func NewVillagerManager() *VillagerManager {
	vm := &VillagerManager{
		villagers: make(map[string]*VillagerType),
	}

	// Initialize default villager types
	vm.villagers["villager"] = &VillagerType{
		Count:    0,
		FoodCost: 5,
		Assignment: VillagerAssignment{
			"food":      0,
			"wood":      0,
			"stone":     0,
			"gold":      0,
			"knowledge": 0,
			"idle":      0,
		},
	}

	vm.villagers["scholar"] = &VillagerType{
		Count:    0,
		FoodCost: 8,
		Assignment: VillagerAssignment{
			"knowledge": 0,
			"idle":      0,
		},
	}

	return vm
}

// Add adds new villagers
func (vm *VillagerManager) Add(villagerType string, count int) bool {
	if v, exists := vm.villagers[villagerType]; exists {
		v.Count += count
		v.Assignment["idle"] += count
		return true
	}
	return false
}

// Remove villagers
func (vm *VillagerManager) Remove(villagerType string, count int) bool {
	if v, exists := vm.villagers[villagerType]; exists && v.Count >= count {
		v.Count -= count

		// Remove from idle if possible
		idleCount := v.Assignment["idle"]
		if idleCount >= count {
			v.Assignment["idle"] -= count
		} else {
			// Need to remove from assignments
			v.Assignment["idle"] = 0
			remaining := count - idleCount

			// Remove remaining villagers from assignments
			for resource, assigned := range v.Assignment {
				if resource == "idle" {
					continue
				}

				if assigned > 0 {
					if assigned >= remaining {
						v.Assignment[resource] -= remaining
						remaining = 0
						break
					} else {
						remaining -= assigned
						v.Assignment[resource] = 0
					}
				}
			}
		}
		return true
	}
	return false
}

// Assigns villagers to gather a resource
func (vm *VillagerManager) Assign(villagerType, resource string, count int) bool {
	if v, exists := vm.villagers[villagerType]; exists {
		// Check if this villager type can gather this resource
		if _, canGather := v.Assignment[resource]; !canGather {
			return false
		}

		// Check if we have enough idle villagers
		if v.Assignment["idle"] < count {
			return false
		}

		// Assign villagers
		v.Assignment[resource] += count
		v.Assignment["idle"] -= count
		return true
	}
	return false
}

// Unassign villagers from a resource
func (vm *VillagerManager) Unassign(villagerType, resource string, count int) bool {
	if v, exists := vm.villagers[villagerType]; exists {
		// Check if this villager type can gather this resource
		if _, canGather := v.Assignment[resource]; !canGather {
			return false
		}

		// Check if we have enough assigned villagers
		if v.Assignment[resource] < count {
			return false
		}

		// Unassign villagers
		v.Assignment[resource] -= count
		v.Assignment["idle"] += count
		return true
	}
	return false
}

// GetCount returns the count of a specific villager type
func (vm *VillagerManager) GetCount(villagerType string) int {
	if v, exists := vm.villagers[villagerType]; exists {
		return v.Count
	}
	return 0
}

// GetAll returns all villagers and their info
type VillagerInfo struct {
	Count      int
	Assignment map[string]int
}

func (vm *VillagerManager) GetAll() map[string]VillagerInfo {
	result := make(map[string]VillagerInfo)
	for vtype, v := range vm.villagers {
		result[vtype] = VillagerInfo{
			Count:      v.Count,
			Assignment: v.Assignment,
		}
	}
	return result
}

// GetFoodConsumption calculates total food consumption
func (vm *VillagerManager) GetFoodConsumption() float64 {
	var total float64
	for _, v := range vm.villagers {
		total += float64(v.Count) * v.FoodCost
	}
	return total
}

// CollectResources collects resources based on villager assignments
func (vm *VillagerManager) CollectResources(rm *ResourceManager) {
	for _, v := range vm.villagers {
		for resource, count := range v.Assignment {
			if resource != "idle" && count > 0 {
				// Calculate collection amount
				collectionRate := rm.GetCollectionRate(resource)
				amount := float64(count) * collectionRate

				// Add the resources
				rm.Add(resource, amount)
			}
		}
	}

	// Consume food
	foodConsumption := vm.GetFoodConsumption()
	rm.Remove("food", foodConsumption)
}

// CollectResourcesAndTrack collects resources based on villager assignments and tracks statistics
func (vm *VillagerManager) CollectResourcesAndTrack(rm *ResourceManager, stats *GameStats) {
	for _, v := range vm.villagers {
		for resource, count := range v.Assignment {
			if resource != "idle" && count > 0 {
				// Calculate collection amount
				collectionRate := rm.GetCollectionRate(resource)
				amount := float64(count) * collectionRate

				// Add the resources
				rm.Add(resource, amount)

				// Track resource gathering in stats
				stats.AddResourceGathered(resource, amount)
			}
		}
	}

	// Consume food
	foodConsumption := vm.GetFoodConsumption()
	rm.Remove("food", foodConsumption)
}
