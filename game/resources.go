package game

// ResourceManager handles all game resources
type ResourceManager struct {
	resources       map[string]float64
	collectionRates map[string]float64
}

// NewResourceManager creates a new resource manager
func NewResourceManager() *ResourceManager {
	rm := &ResourceManager{
		resources: map[string]float64{
			"food":      0,
			"wood":      0,
			"stone":     0,
			"gold":      0,
			"knowledge": 0,
		},
		collectionRates: map[string]float64{
			"food":      1.0,
			"wood":      1.0,
			"stone":     0.5,
			"gold":      0.2,
			"knowledge": 0.1,
		},
	}
	return rm
}

// Add adds resources to the inventory
func (rm *ResourceManager) Add(resource string, amount float64) bool {
	if _, exists := rm.resources[resource]; exists {
		rm.resources[resource] += amount
		return true
	}
	return false
}

// Remove removes resources from the inventory
func (rm *ResourceManager) Remove(resource string, amount float64) bool {
	if _, exists := rm.resources[resource]; exists && rm.resources[resource] >= amount {
		rm.resources[resource] -= amount
		return true
	}
	return false
}

// Has checks if we have enough of a resource
func (rm *ResourceManager) Has(resource string, amount float64) bool {
	return rm.resources[resource] >= amount
}

// Get returns the amount of a specific resource
func (rm *ResourceManager) Get(resource string) float64 {
	return rm.resources[resource]
}

// GetAll returns all resources and their amounts
func (rm *ResourceManager) GetAll() map[string]float64 {
	return rm.resources
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
