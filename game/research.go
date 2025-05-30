package game

// ResearchManager handles technology research and unlocking new abilities
type ResearchManager struct {
	technologies     map[string]Technology
	researchedTechs  map[string]bool
	currentResearch  string
	researchProgress float64
}

// Technology represents a researchable technology
type Technology struct {
	Name        string
	Description string
	Age         string
	Cost        float64
	Prerequisites []string
	Unlocks     map[string]interface{}
}

// NewResearchManager creates a new research manager
func NewResearchManager() *ResearchManager {
	rm := &ResearchManager{
		technologies:    make(map[string]Technology),
		researchedTechs: make(map[string]bool),
		currentResearch: "",
	}

	// Define technologies
	rm.technologies["agriculture"] = Technology{
		Name:        "Agriculture",
		Description: "Improve food production methods",
		Age:         "Stone Age",
		Cost:        20,
		Prerequisites: []string{},
		Unlocks: map[string]interface{}{
			"food_production_bonus": 0.2,
		},
	}

	rm.technologies["toolmaking"] = Technology{
		Name:        "Toolmaking",
		Description: "Develop better tools for resource gathering",
		Age:         "Stone Age",
		Cost:        25,
		Prerequisites: []string{},
		Unlocks: map[string]interface{}{
			"resource_production_bonus": 0.1,
		},
	}

	rm.technologies["writing"] = Technology{
		Name:        "Writing",
		Description: "Develop a writing system to record knowledge",
		Age:         "Bronze Age",
		Cost:        40,
		Prerequisites: []string{},
		Unlocks: map[string]interface{}{
			"knowledge_production_bonus": 0.2,
		},
	}

	rm.technologies["metallurgy"] = Technology{
		Name:        "Metallurgy",
		Description: "Learn how to work with metals",
		Age:         "Bronze Age",
		Cost:        50,
		Prerequisites: []string{},
		Unlocks: map[string]interface{}{
			"new_building": "foundry",
		},
	}

	rm.technologies["mathematics"] = Technology{
		Name:        "Mathematics",
		Description: "Develop mathematical concepts",
		Age:         "Iron Age",
		Cost:        60,
		Prerequisites: []string{"writing"},
		Unlocks: map[string]interface{}{
			"knowledge_production_bonus": 0.3,
			"resource_production_bonus": 0.1,
		},
	}

	return rm
}

// StartResearch begins research on a technology
func (rm *ResearchManager) StartResearch(techName string, knowledgePoints float64) bool {
	// Check if technology exists
	tech, exists := rm.technologies[techName]
	if !exists {
		return false
	}

	// Check if already researched
	if rm.researchedTechs[techName] {
		return false
	}

	// Check prerequisites
	for _, prereq := range tech.Prerequisites {
		if !rm.researchedTechs[prereq] {
			return false
		}
	}

	// Start research
	rm.currentResearch = techName
	rm.researchProgress = knowledgePoints
	return true
}

// ContinueResearch adds knowledge points to current research
func (rm *ResearchManager) ContinueResearch(knowledgePoints float64) (string, bool) {
	if rm.currentResearch == "" {
		return "", false
	}

	tech := rm.technologies[rm.currentResearch]
	rm.researchProgress += knowledgePoints

	// Check if research is complete
	if rm.researchProgress >= tech.Cost {
		completedTech := rm.currentResearch
		rm.researchedTechs[completedTech] = true
		rm.currentResearch = ""
		rm.researchProgress = 0
		return completedTech, true
	}

	return "", false
}

// GetProgress returns the current research progress
func (rm *ResearchManager) GetProgress() (string, float64, float64) {
	if rm.currentResearch == "" {
		return "", 0, 0
	}

	tech := rm.technologies[rm.currentResearch]
	return rm.currentResearch, rm.researchProgress, tech.Cost
}

// IsResearched checks if a technology has been researched
func (rm *ResearchManager) IsResearched(techName string) bool {
	return rm.researchedTechs[techName]
}

// GetAvailableTechnologies returns technologies available for research in the current age
func (rm *ResearchManager) GetAvailableTechnologies(currentAge string) map[string]Technology {
	result := make(map[string]Technology)
	
	// Get age index
	ageIndex := 0
	ages := []string{"Stone Age", "Bronze Age", "Iron Age", "Medieval Age", "Renaissance Age", "Industrial Age", "Modern Age"}
	for i, age := range ages {
		if age == currentAge {
			ageIndex = i
			break
		}
	}
	
	// Get technologies from current and previous ages
	for name, tech := range rm.technologies {
		techAgeIndex := 0
		for i, age := range ages {
			if age == tech.Age {
				techAgeIndex = i
				break
			}
		}
		
		// Technology is from current or previous age and not already researched
		if techAgeIndex <= ageIndex && !rm.researchedTechs[name] {
			// Check prerequisites
			allPrereqsMet := true
			for _, prereq := range tech.Prerequisites {
				if !rm.researchedTechs[prereq] {
					allPrereqsMet = false
					break
				}
			}
			
			if allPrereqsMet {
				result[name] = tech
			}
		}
	}
	
	return result
}

// GetResearchedTechnologies returns all researched technologies
func (rm *ResearchManager) GetResearchedTechnologies() map[string]Technology {
	result := make(map[string]Technology)
	for name := range rm.researchedTechs {
		if tech, exists := rm.technologies[name]; exists && rm.researchedTechs[name] {
			result[name] = tech
		}
	}
	return result
}

// GetAllTechnologies returns all technologies
func (rm *ResearchManager) GetAllTechnologies() map[string]Technology {
	return rm.technologies
}

// ApplyResearchBonuses applies research bonuses to resource production
func (rm *ResearchManager) ApplyResearchBonuses(resourceRates map[string]float64) map[string]float64 {
	result := make(map[string]float64)
	for resource, rate := range resourceRates {
		result[resource] = rate
	}
	
	// Apply global resource production bonus
	globalBonus := 0.0
	for tech, researched := range rm.researchedTechs {
		if researched {
			if bonus, exists := rm.technologies[tech].Unlocks["resource_production_bonus"]; exists {
				globalBonus += bonus.(float64)
			}
		}
	}
	
	if globalBonus > 0 {
		for resource, rate := range result {
			result[resource] = rate * (1 + globalBonus)
		}
	}
	
	// Apply specific resource bonuses
	for tech, researched := range rm.researchedTechs {
		if researched {
			for unlock, value := range rm.technologies[tech].Unlocks {
				if resource := getResourceBonusType(unlock); resource != "" {
					if bonus, ok := value.(float64); ok {
						if _, exists := result[resource]; exists {
							result[resource] = result[resource] * (1 + bonus)
						}
					}
				}
			}
		}
	}
	
	return result
}

// Helper function to extract resource type from bonus name
func getResourceBonusType(bonusName string) string {
	resourceTypes := []string{"food", "wood", "stone", "gold", "knowledge"}
	for _, resource := range resourceTypes {
		if bonusName == resource+"_production_bonus" {
			return resource
		}
	}
	return ""
}
