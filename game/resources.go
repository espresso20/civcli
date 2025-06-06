package game

// ResourceManager handles all game resources
type ResourceManager struct {
	resources       map[string]float64
	collectionRates map[string]float64
	foodSources     []string // List of resources that count as food sources
}

// NewResourceManager creates a new resource manager
func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		resources: map[string]float64{
			"foraging":  0,
			"wood":      0,
			"stone":     0,
			"gold":      0,
			"knowledge": 0,
			"hunting":   0,
		},
		collectionRates: map[string]float64{
			"foraging":  1.0,
			"wood":      1.0,
			"stone":     0.5,
			"gold":      0.2,
			"knowledge": 0.1,
			"hunting":   1.8,
		},
		foodSources: []string{"foraging", "hunting"}, // Define which resources count as food
	}
	return rm
}

// Add adds resources to the inventory
func (rm *ResourceManager) Add(resource string, amount float64) bool {
	// Special case for "food" - add to "foraging" instead
	if resource == "food" {
		if _, exists := rm.resources["foraging"]; exists {
			rm.resources["foraging"] += amount
			return true
		}
		return false
	}

	// Normal case - add to specific resource
	if _, exists := rm.resources[resource]; exists {
		rm.resources[resource] += amount
		return true
	}
	return false
}

// Remove removes resources from the inventory
func (rm *ResourceManager) Remove(resource string, amount float64) bool {
	// Special case for "food" - delegate to RemoveFood method
	if resource == "food" {
		return rm.RemoveFood(amount)
	}

	// Normal case - remove from specific resource
	if _, exists := rm.resources[resource]; exists && rm.resources[resource] >= amount {
		rm.resources[resource] -= amount
		return true
	}
	return false
}

// checks if we have enough of a resource
func (rm *ResourceManager) Has(resource string, amount float64) bool {
	// Special case for "food" - delegate to HasFood method
	if resource == "food" {
		return rm.HasFood(amount)
	}

	// Normal case - check specific resource
	if val, exists := rm.resources[resource]; exists {
		return val >= amount
	}
	return false
}

// Get returns the amount of a specific resource
func (rm *ResourceManager) Get(resource string) float64 {
	return rm.resources[resource]
}

// GetAll returns all resources and their amounts
func (rm *ResourceManager) GetAll() map[string]float64 {
	// Create a copy of the resources map
	resources := make(map[string]float64)
	for key, value := range rm.resources {
		resources[key] = value
	}

	// We don't add the virtual "food" resource to avoid confusion in the UI
	// Food is managed internally but not shown as a separate resource

	return resources
}

// GetCollectionRate returns the collection rate for a resource
func (rm *ResourceManager) GetCollectionRate(resource string) float64 {
	return rm.collectionRates[resource]
}

// SetCollectionRate sets the collection rate for a resource
func (rm *ResourceManager) SetCollectionRate(resource string, rate float64) bool {
	if _, exists := rm.collectionRates[resource]; exists {
		rm.collectionRates[resource] = rate
		return true
	}
	return false
}

// GetTotalFood returns the sum of all food resources
func (rm *ResourceManager) GetTotalFood() float64 {
	var total float64
	for _, foodSource := range rm.foodSources {
		total += rm.resources[foodSource]
	}
	return total
}

// HasFood checks if there is enough total food available
func (rm *ResourceManager) HasFood(amount float64) bool {
	return rm.GetTotalFood() >= amount
}

// RemoveFood removes food proportionally from all food sources
func (rm *ResourceManager) RemoveFood(amount float64) bool {
	if !rm.HasFood(amount) {
		return false
	}

	totalFood := rm.GetTotalFood()

	// Nothing to remove or negative amount
	if totalFood <= 0 || amount <= 0 {
		return true
	}

	// Remove proportionally from each food source
	for _, foodSource := range rm.foodSources {
		if rm.resources[foodSource] > 0 {
			// Calculate proportion of this food source
			proportion := rm.resources[foodSource] / totalFood
			// Remove proportional amount
			amountToRemove := amount * proportion
			rm.resources[foodSource] -= amountToRemove

			// Handle floating point precision issues
			if rm.resources[foodSource] < 0.00001 {
				rm.resources[foodSource] = 0
			}
		}
	}

	return true
}
