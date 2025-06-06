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
		FoodCost: 0.5,
		Assignment: VillagerAssignment{
			"foraging":  0,
			"wood":      0,
			"stone":     0,
			"gold":      0,
			"knowledge": 0,
			"hunting":   0,
			"idle":      0,
		},
	}

	vm.villagers["scholar"] = &VillagerType{
		Count:    0,
		FoodCost: 0.75,
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
func (vm *VillagerManager) CollectResources(rm *ResourceManager, bm *BuildingManager) {
	// First, collect all resources by villager type and assignment
	vm.gatherAllResources(rm, bm)

	// Then consume food from the total food pool
	foodConsumption := vm.GetFoodConsumption()
	rm.RemoveFood(foodConsumption)
}

// gatherAllResources handles the resource gathering for all villager types
func (vm *VillagerManager) gatherAllResources(rm *ResourceManager, bm *BuildingManager) {
	for vtype, v := range vm.villagers {
		for resource, count := range v.Assignment {
			if resource == "idle" || count <= 0 {
				continue
			}

			// Get the base collection rate for this resource
			baseRate := rm.GetCollectionRate(resource)

			// Calculate the final amount based on resource type
			switch resource {
			case "knowledge":
				vm.gatherKnowledge(rm, bm, vtype, count, baseRate)
			case "hunting":
				vm.gatherHunting(rm, bm, vtype, count, baseRate)
			default:
				vm.gatherStandardResource(rm, bm, vtype, resource, count, baseRate)
			}
		}
	}
}

// gatherKnowledge handles specialized knowledge gathering with different rates by villager type
func (vm *VillagerManager) gatherKnowledge(rm *ResourceManager, bm *BuildingManager, vtype string, count int, baseRate float64) {
	// Apply villager-specific knowledge gathering modifiers
	modifiedRate := baseRate
	if vtype == "villager" {
		// Regular villagers gather knowledge at 20% of the normal rate
		modifiedRate *= 0.2
	} else if vtype == "scholar" {
		// Scholars gather knowledge at 150% of the normal rate
		modifiedRate *= 1.5
	}

	// Apply building bonuses
	buildingBonus := bm.GetCollectionRateBonus(vtype, "knowledge")
	modifiedRate *= (1.0 + buildingBonus)

	// Calculate final amount and add the resource
	amount := float64(count) * modifiedRate
	rm.Add("knowledge", amount)
}

// gatherHunting handles hunting which provides both hunting resource and food bonus
func (vm *VillagerManager) gatherHunting(rm *ResourceManager, bm *BuildingManager, vtype string, count int, baseRate float64) {
	// Apply building bonuses to hunting rate
	buildingBonus := bm.GetCollectionRateBonus(vtype, "hunting")
	modifiedRate := baseRate * (1.0 + buildingBonus)

	// Calculate hunting resource amount
	amount := float64(count) * modifiedRate
	rm.Add("hunting", amount)

	// Add bonus food from hunting (40% of hunting collection)
	foodBonus := baseRate * 0.4 * float64(count)
	rm.Add("food", foodBonus)
}

// gatherStandardResource handles standard resource gathering without special rules
func (vm *VillagerManager) gatherStandardResource(rm *ResourceManager, bm *BuildingManager, vtype string, resource string, count int, baseRate float64) {
	// Apply building bonuses
	buildingBonus := bm.GetCollectionRateBonus(vtype, resource)
	modifiedRate := baseRate * (1.0 + buildingBonus)

	// Calculate final amount and add the resource
	amount := float64(count) * modifiedRate
	rm.Add(resource, amount)
}

// CollectResourcesAndTrack collects resources based on villager assignments and tracks statistics
func (vm *VillagerManager) CollectResourcesAndTrack(rm *ResourceManager, stats *GameStats, bm *BuildingManager) {
	// Use the refactored resource gathering approach while tracking statistics
	vm.gatherAllResourcesAndTrack(rm, bm, stats)

	// Then consume food from the total food pool
	foodConsumption := vm.GetFoodConsumption()
	rm.RemoveFood(foodConsumption)
}

// gatherAllResourcesAndTrack handles resource gathering with statistics tracking
func (vm *VillagerManager) gatherAllResourcesAndTrack(rm *ResourceManager, bm *BuildingManager, stats *GameStats) {
	for vtype, v := range vm.villagers {
		for resource, count := range v.Assignment {
			if resource == "idle" || count <= 0 {
				continue
			}

			// Get the base collection rate for this resource
			baseRate := rm.GetCollectionRate(resource)

			// Calculate the final amount based on resource type and track statistics
			switch resource {
			case "knowledge":
				amount := vm.gatherKnowledgeWithTracking(rm, bm, vtype, count, baseRate, stats)
				stats.AddResourceGathered(resource, amount)
			case "hunting":
				huntingAmount, foodAmount := vm.gatherHuntingWithTracking(rm, bm, vtype, count, baseRate)
				stats.AddResourceGathered(resource, huntingAmount)
				stats.AddResourceGathered("food", foodAmount)
			default:
				amount := vm.gatherStandardResourceWithTracking(rm, bm, vtype, resource, count, baseRate)
				stats.AddResourceGathered(resource, amount)
			}
		}
	}
}

// gatherKnowledgeWithTracking handles specialized knowledge gathering with tracking
func (vm *VillagerManager) gatherKnowledgeWithTracking(rm *ResourceManager, bm *BuildingManager, vtype string, count int, baseRate float64, stats *GameStats) float64 {
	// Apply villager-specific knowledge gathering modifiers
	modifiedRate := baseRate
	if vtype == "villager" {
		// Regular villagers gather knowledge at 20% of the normal rate
		modifiedRate *= 0.2
	} else if vtype == "scholar" {
		// Scholars gather knowledge at 150% of the normal rate
		modifiedRate *= 1.5
	}

	// Apply building bonuses
	buildingBonus := bm.GetCollectionRateBonus(vtype, "knowledge")
	modifiedRate *= (1.0 + buildingBonus)

	// Calculate final amount and add the resource
	amount := float64(count) * modifiedRate
	rm.Add("knowledge", amount)

	return amount
}

// gatherHuntingWithTracking handles hunting which provides both hunting resource and food bonus with tracking
func (vm *VillagerManager) gatherHuntingWithTracking(rm *ResourceManager, bm *BuildingManager, vtype string, count int, baseRate float64) (float64, float64) {
	// Apply building bonuses to hunting rate
	buildingBonus := bm.GetCollectionRateBonus(vtype, "hunting")
	modifiedRate := baseRate * (1.0 + buildingBonus)

	// Calculate hunting resource amount
	huntingAmount := float64(count) * modifiedRate
	rm.Add("hunting", huntingAmount)

	// Add bonus food from hunting (40% of hunting collection)
	foodBonus := baseRate * 0.4 * float64(count)
	rm.Add("food", foodBonus)

	return huntingAmount, foodBonus
}

// gatherStandardResourceWithTracking handles standard resource gathering with tracking
func (vm *VillagerManager) gatherStandardResourceWithTracking(rm *ResourceManager, bm *BuildingManager, vtype string, resource string, count int, baseRate float64) float64 {
	// Apply building bonuses
	buildingBonus := bm.GetCollectionRateBonus(vtype, resource)
	modifiedRate := baseRate * (1.0 + buildingBonus)

	// Calculate final amount and add the resource
	amount := float64(count) * modifiedRate
	rm.Add(resource, amount)

	return amount
}
