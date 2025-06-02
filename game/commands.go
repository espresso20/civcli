package game

import (
	"strconv"
	"strings"
)

// CommandHandler processes user commands
type CommandHandler struct {
	Game     *GameEngine
	Commands map[string]string
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(gameEngine *GameEngine) *CommandHandler {
	ch := &CommandHandler{
		Game: gameEngine,
		Commands: map[string]string{
			"help":      "Display available commands",
			"gather":    "Assign villagers to gather resources (gather <resource> <count>)",
			"build":     "Build a structure (build <building>)",
			"status":    "Show detailed status of your civilization",
			"assign":    "Assign villagers to tasks (assign <villager_type> <resource> <count>)",
			"unassign":  "Unassign villagers from tasks (unassign <villager_type> <resource> <count>)",
			"recruit":   "Recruit new villagers (recruit <villager_type> <count>)",
			"buildings": "List available buildings and their costs",
			"research":  "Start researching a technology (research <technology>)",
			"techs":     "List available technologies for research",
			"library":   "Access the in-game library (library [topic])",
			"save":      "Save the current game (save <filename>)",
			"load":      "Load a saved game (load <filename>)",
			"saves":     "List all saved games",
			"stats":     "Display game statistics",
			"clear":     "Clear the console screen",
			"quit":      "Exit the game",
		},
	}
	return ch
}

// GetCommandList returns the list of available commands
func (ch *CommandHandler) GetCommandList() []string {
	commands := make([]string, 0, len(ch.Commands))
	for cmd := range ch.Commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Process processes a command string
func (ch *CommandHandler) Process(commandStr string) {
	// Empty command
	if commandStr == "" {
		return
	}

	// Split the command into parts
	parts := strings.Fields(commandStr)
	command := strings.ToLower(parts[0])
	args := parts[1:]

	// Process command
	switch command {
	case "help":
		ch.CmdHelp()
	case "gather":
		ch.CmdGather(args)
	case "build":
		ch.CmdBuild(args)
	case "status":
		ch.CmdStatus()
	case "assign":
		ch.CmdAssign(args)
	case "unassign":
		ch.CmdUnassign(args)
	case "recruit":
		ch.CmdRecruit(args)
	case "buildings":
		ch.CmdBuildings()
	case "research":
		ch.CmdResearch(args)
	case "techs":
		ch.CmdTechs()
	case "save":
		ch.CmdSave(args)
	case "load":
		ch.CmdLoad(args)
	case "saves":
		ch.CmdListSaves()
	case "stats":
		ch.CmdStats()
	case "library":
		ch.CmdLibrary(args)
	case "clear":
		// This will be handled in the UI
	case "quit":
		ch.Game.Quit()
	default:
		ch.Game.Display.ShowMessage("Unknown command: "+command+". Type 'help' for available commands.", "error")
	}
}

// CmdHelp displays help information
func (ch *CommandHandler) CmdHelp() {
	ch.Game.Display.ShowHelp(ch.Commands)
}

// CmdGather assigns villagers to gather resources
func (ch *CommandHandler) CmdGather(args []string) {
	if len(args) != 2 {
		ch.Game.Display.ShowMessage("Usage: gather <resource> <count>", "error")
		return
	}

	resource := args[0]
	count, err := strconv.Atoi(args[1])
	if err != nil || count <= 0 {
		ch.Game.Display.ShowMessage("Count must be a positive number", "error")
		return
	}

	// Try to assign villagers
	if ch.Game.Villagers.Assign("villager", resource, count) {
		ch.Game.Display.ShowMessage("Assigned "+strconv.Itoa(count)+" villagers to gather "+resource, "success")
	} else {
		ch.Game.Display.ShowMessage("Failed to assign villagers. Not enough idle villagers or invalid resource.", "error")
	}
}

// CmdBuild builds a structure
func (ch *CommandHandler) CmdBuild(args []string) {
	if len(args) != 1 {
		ch.Game.Display.ShowMessage("Usage: build <building>", "error")
		return
	}

	building := args[0]

	// Check if building is available in current age
	currentAgeIndex := ch.Game.Progress.GetCurrentAgeIndex(ch.Game.Age)
	availableBuildings := []string{}

	for i, age := range ch.Game.Progress.GetAllAges() {
		if i <= currentAgeIndex {
			availableBuildings = append(availableBuildings, ch.Game.Progress.GetUnlocks(age).Buildings...)
		}
	}

	buildingAvailable := false
	for _, b := range availableBuildings {
		if b == building {
			buildingAvailable = true
			break
		}
	}

	if !buildingAvailable {
		ch.Game.Display.ShowMessage(building+" is not available in the "+ch.Game.Age, "error")
		return
	}

	// Try to build
	if ch.Game.Buildings.Build(building, ch.Game.Resources) {
		ch.Game.Display.ShowMessage("Built a new "+building, "success")

		// Track building in stats
		ch.Game.Stats.AddEvent(ch.Game.Tick, "building_built", "Built a new "+building)
		ch.Game.Stats.AddBuildingBuilt(building)
	} else {
		costs := ch.Game.Buildings.GetCost(building)
		costStrs := []string{}
		for res, amount := range costs {
			costStrs = append(costStrs, strconv.FormatFloat(amount, 'f', 0, 64)+" "+res)
		}
		costStr := strings.Join(costStrs, ", ")
		ch.Game.Display.ShowMessage("Failed to build "+building+". Required resources: "+costStr, "error")
	}
}

// CmdStatus shows detailed status
func (ch *CommandHandler) CmdStatus() {
	// This is handled by the UI
}

// CmdAssign assigns villagers to tasks
func (ch *CommandHandler) CmdAssign(args []string) {
	if len(args) != 3 {
		ch.Game.Display.ShowMessage("Usage: assign <villager_type> <resource> <count>", "error")
		return
	}

	villagerType := args[0]
	resource := args[1]
	count, err := strconv.Atoi(args[2])
	if err != nil || count <= 0 {
		ch.Game.Display.ShowMessage("Count must be a positive number", "error")
		return
	}

	// Try to assign villagers
	if ch.Game.Villagers.Assign(villagerType, resource, count) {
		ch.Game.Display.ShowMessage("Assigned "+strconv.Itoa(count)+" "+villagerType+"s to "+resource, "success")
	} else {
		ch.Game.Display.ShowMessage("Failed to assign "+villagerType+"s. Not enough idle villagers or invalid resource.", "error")
	}
}

// CmdUnassign unassigns villagers from tasks
func (ch *CommandHandler) CmdUnassign(args []string) {
	if len(args) != 3 {
		ch.Game.Display.ShowMessage("Usage: unassign <villager_type> <resource> <count>", "error")
		return
	}

	villagerType := args[0]
	resource := args[1]
	count, err := strconv.Atoi(args[2])
	if err != nil || count <= 0 {
		ch.Game.Display.ShowMessage("Count must be a positive number", "error")
		return
	}

	// Try to unassign villagers
	if ch.Game.Villagers.Unassign(villagerType, resource, count) {
		ch.Game.Display.ShowMessage("Unassigned "+strconv.Itoa(count)+" "+villagerType+"s from "+resource, "success")
	} else {
		ch.Game.Display.ShowMessage("Failed to unassign "+villagerType+"s. Not enough assigned to "+resource+".", "error")
	}
}

// CmdRecruit recruits new villagers
func (ch *CommandHandler) CmdRecruit(args []string) {
	if len(args) != 2 {
		ch.Game.Display.ShowMessage("Usage: recruit <villager_type> <count>", "error")
		return
	}

	villagerType := args[0]
	count, err := strconv.Atoi(args[1])
	if err != nil || count <= 0 {
		ch.Game.Display.ShowMessage("Count must be a positive number", "error")
		return
	}

	// Check if villager type is available in current age
	currentAgeIndex := ch.Game.Progress.GetCurrentAgeIndex(ch.Game.Age)
	availableVillagers := []string{"villager"} // Villager always available

	for i, age := range ch.Game.Progress.GetAllAges() {
		if i <= currentAgeIndex {
			availableVillagers = append(availableVillagers, ch.Game.Progress.GetUnlocks(age).Villagers...)
		}
	}

	villagerAvailable := false
	for _, v := range availableVillagers {
		if v == villagerType {
			villagerAvailable = true
			break
		}
	}

	if !villagerAvailable {
		ch.Game.Display.ShowMessage(villagerType+" is not available in the "+ch.Game.Age, "error")
		return
	}

	// Check villager capacity
	capacity := ch.Game.Buildings.GetVillagerCapacity()
	totalVillagers := 0
	for _, v := range ch.Game.Villagers.GetAll() {
		totalVillagers += v.Count
	}

	if totalVillagers+count > capacity {
		ch.Game.Display.ShowMessage("Not enough housing capacity. Current: "+strconv.Itoa(totalVillagers)+"/"+strconv.Itoa(capacity), "error")
		return
	}

	// Get food cost
	foodCost := float64(0)
	for vtype, info := range ch.Game.Villagers.villagers {
		if vtype == villagerType {
			foodCost = info.FoodCost * float64(count)
			break
		}
	}

	// Check food
	if !ch.Game.Resources.Has("food", foodCost) {
		ch.Game.Display.ShowMessage("Not enough food. Need "+strconv.FormatFloat(foodCost, 'f', 0, 64)+" food.", "error")
		return
	}

	// Recruit villagers
	ch.Game.Resources.Remove("food", foodCost)
	ch.Game.Villagers.Add(villagerType, count)
	ch.Game.Display.ShowMessage("Recruited "+strconv.Itoa(count)+" new "+villagerType+"s", "success")

	// Track recruitment in stats
	ch.Game.Stats.AddEvent(ch.Game.Tick, "villager_recruited", "Recruited "+strconv.Itoa(count)+" new "+villagerType+"s")
	for i := 0; i < count; i++ {
		ch.Game.Stats.AddVillagerRecruited(villagerType)
	}
}

// CmdBuildings lists available buildings and their costs
func (ch *CommandHandler) CmdBuildings() {
	// This will be handled by the UI
}

// CmdSave saves the current game
func (ch *CommandHandler) CmdSave(args []string) {
	if len(args) != 1 {
		ch.Game.Display.ShowMessage("Usage: save <filename>", "error")
		return
	}

	filename := args[0]
	err := ch.Game.SaveGame(filename)
	if err != nil {
		ch.Game.Display.ShowMessage("Failed to save game: "+err.Error(), "error")
		return
	}

	ch.Game.Display.ShowMessage("Game saved as '"+filename+"'", "success")
}

// CmdLoad loads a saved game
func (ch *CommandHandler) CmdLoad(args []string) {
	if len(args) != 1 {
		ch.Game.Display.ShowMessage("Usage: load <filename>", "error")
		return
	}

	filename := args[0]
	err := ch.Game.LoadGame(filename)
	if err != nil {
		ch.Game.Display.ShowMessage("Failed to load game: "+err.Error(), "error")
		return
	}

	ch.Game.Display.ShowMessage("Game '"+filename+"' loaded successfully", "success")
}

// CmdListSaves lists all saved games
func (ch *CommandHandler) CmdListSaves() {
	saves, err := ListSaves()
	if err != nil {
		ch.Game.Display.ShowMessage("Failed to list saves: "+err.Error(), "error")
		return
	}

	if len(saves) == 0 {
		ch.Game.Display.ShowMessage("No saved games found", "info")
		return
	}

	ch.Game.Display.ShowMessage("Available saved games:", "info")
	for _, save := range saves {
		ch.Game.Display.ShowMessage("- "+save, "info")
	}
}

// CmdStats shows game statistics
func (ch *CommandHandler) CmdStats() {
	// Calculate play time
	playTime := ch.Game.Stats.GetPlayTime()

	// Show general stats
	ch.Game.Display.ShowMessage("=== Game Statistics ===", "highlight")
	ch.Game.Display.ShowMessage("Play time: "+playTime, "info")
	ch.Game.Display.ShowMessage("Current age: "+ch.Game.Age, "info")
	ch.Game.Display.ShowMessage("Game ticks: "+strconv.Itoa(ch.Game.Tick), "info")

	// Show resources gathered
	ch.Game.Display.ShowMessage("\n=== Resources Gathered ===", "highlight")
	for resource, amount := range ch.Game.Stats.ResourcesGathered {
		ch.Game.Display.ShowMessage(resource+": "+strconv.FormatFloat(amount, 'f', 1, 64), "info")
	}
	totalResources := ch.Game.Stats.GetTotalResourcesGathered()
	ch.Game.Display.ShowMessage("Total resources: "+strconv.FormatFloat(totalResources, 'f', 1, 64), "success")

	// Show buildings built
	ch.Game.Display.ShowMessage("\n=== Buildings Built ===", "highlight")
	for building, count := range ch.Game.Stats.BuildingsBuilt {
		ch.Game.Display.ShowMessage(building+": "+strconv.Itoa(count), "info")
	}
	totalBuildings := ch.Game.Stats.GetTotalBuildingsBuilt()
	ch.Game.Display.ShowMessage("Total buildings: "+strconv.Itoa(totalBuildings), "success")

	// Show villagers recruited
	ch.Game.Display.ShowMessage("\n=== Villagers Recruited ===", "highlight")
	for villager, count := range ch.Game.Stats.VillagersRecruited {
		ch.Game.Display.ShowMessage(villager+": "+strconv.Itoa(count), "info")
	}
	totalVillagers := ch.Game.Stats.GetTotalVillagersRecruited()
	ch.Game.Display.ShowMessage("Total villagers: "+strconv.Itoa(totalVillagers), "success")

	// Show ages reached
	ch.Game.Display.ShowMessage("\n=== Ages Reached ===", "highlight")
	for _, age := range ch.Game.Stats.AgesReached {
		ch.Game.Display.ShowMessage(age, "info")
	}

	// Show recent events
	ch.Game.Display.ShowMessage("\n=== Recent Events ===", "highlight")
	eventCount := len(ch.Game.Stats.Events)
	startIdx := 0
	if eventCount > 10 {
		startIdx = eventCount - 10
	}

	for i := startIdx; i < eventCount; i++ {
		event := ch.Game.Stats.Events[i]
		ch.Game.Display.ShowMessage(
			"Tick "+strconv.Itoa(event.Tick)+": "+event.Message,
			"info",
		)
	}
}

// CmdResearch starts researching a technology
func (ch *CommandHandler) CmdResearch(args []string) {
	if len(args) != 1 {
		ch.Game.Display.ShowMessage("Usage: research <technology>", "error")
		return
	}

	techName := args[0]

	// Check if the technology is available
	availableTechs := ch.Game.Research.GetAvailableTechnologies(ch.Game.Age)
	tech, exists := availableTechs[techName]

	if !exists {
		ch.Game.Display.ShowMessage("Technology '"+techName+"' is not available for research.", "error")
		return
	}

	// Check if we're already researching something
	currentTech, progress, _ := ch.Game.Research.GetProgress()
	if currentTech != "" {
		ch.Game.Display.ShowMessage("You are already researching "+currentTech+" ("+
			strconv.FormatFloat(progress, 'f', 1, 64)+" / "+
			strconv.FormatFloat(tech.Cost, 'f', 1, 64)+")", "error")
		return
	}

	// Check if the player has any knowledge points
	knowledgePoints := ch.Game.Resources.Get("knowledge")
	if knowledgePoints <= 0 {
		ch.Game.Display.ShowMessage("You cannot research any technology without knowledge points. Assign villagers to gather knowledge.", "error")
		return
	}

	// Start research
	if ch.Game.Research.StartResearch(techName, 0) {
		ch.Game.Display.ShowMessage("Started researching "+techName, "success")
		ch.Game.Stats.AddEvent(ch.Game.Tick, "research_started", "Started researching "+techName)
	} else {
		ch.Game.Display.ShowMessage("Failed to start research on "+techName, "error")
	}
}

// CmdTechs lists available technologies
func (ch *CommandHandler) CmdTechs() {
	// Get available technologies
	availableTechs := ch.Game.Research.GetAvailableTechnologies(ch.Game.Age)

	if len(availableTechs) == 0 {
		ch.Game.Display.ShowMessage("No technologies available for research in the "+ch.Game.Age, "info")
		return
	}

	// Check if the player has any knowledge points
	knowledgePoints := ch.Game.Resources.Get("knowledge")
	if knowledgePoints <= 0 {
		ch.Game.Display.ShowMessage("Note: You need knowledge points to start researching. Assign villagers to gather knowledge.", "warning")
	}

	// Display available technologies
	ch.Game.Display.ShowMessage("=== Available Technologies ===", "highlight")
	for name, tech := range availableTechs {
		ch.Game.Display.ShowMessage(name+": "+tech.Description+" (Cost: "+
			strconv.FormatFloat(tech.Cost, 'f', 0, 64)+" knowledge)", "info")
	}

	// Display current research if any
	currentTech, progress, cost := ch.Game.Research.GetProgress()
	if currentTech != "" {
		ch.Game.Display.ShowMessage("\n=== Current Research ===", "highlight")
		ch.Game.Display.ShowMessage(currentTech+": "+
			strconv.FormatFloat(progress, 'f', 1, 64)+" / "+
			strconv.FormatFloat(cost, 'f', 1, 64)+" ("+
			strconv.FormatFloat(progress/cost*100, 'f', 1, 64)+"%)", "success")
	}

	// Display researched technologies
	researchedTechs := ch.Game.Research.GetResearchedTechnologies()
	if len(researchedTechs) > 0 {
		ch.Game.Display.ShowMessage("\n=== Researched Technologies ===", "highlight")
		for name, tech := range researchedTechs {
			ch.Game.Display.ShowMessage(name+": "+tech.Description, "info")
		}
	}
}

// CmdLibrary displays the in-game library content
func (ch *CommandHandler) CmdLibrary(args []string) {
	// If no topic specified, show the list of available topics (CLI-only)
	if len(args) == 0 {
		ch.Game.Display.ShowLibraryTopicsList(ch.Game.Library.GetTopicList())
		return
	}

	// Get the requested topic
	topicID := strings.ToLower(args[0])
	topic := ch.Game.Library.GetTopic(topicID)

	if topic == nil {
		// Try to search for the topic
		searchResults := ch.Game.Library.SearchTopics(topicID)

		if len(searchResults) == 0 {
			ch.Game.Display.ShowMessage("Topic not found: "+topicID, "error")
			return
		} else if len(searchResults) == 1 {
			// If only one result, show that topic
			for id := range searchResults {
				topic = ch.Game.Library.GetTopic(id)
				break
			}
		} else {
			// Multiple results, show the list of matching topics (CLI-only)
			ch.Game.Display.ShowLibraryTopicsList(searchResults)
			return
		}
	}

	if topic != nil {
		ch.Game.Display.ShowLibraryContent(topic.Title, topic.Content)
	} else {
		ch.Game.Display.ShowMessage("Topic not found: "+topicID, "error")
	}
}
