package game

import (
	"fmt"
	"strings"
)

// LibraryTopic represents a topic in the game's help library
type LibraryTopic struct {
	Title   string
	Content string
}

// LibrarySystem manages the in-game help library
type LibrarySystem struct {
	topics map[string]*LibraryTopic
}

// NewLibrarySystem creates a new library system with pre-populated topics
func NewLibrarySystem() *LibrarySystem {
	ls := &LibrarySystem{
		topics: make(map[string]*LibraryTopic),
	}

	// Populate the library with topics
	ls.AddTopic("villagers", "Villagers Guide", `
[::b][#3498db]VILLAGERS[#ffffff][::b]

Villagers are the lifeblood of your civilization. They gather resources, 
occupy buildings, and help your civilization grow and advance.

[#f1c40f]Recruiting Villagers:[#ffffff]
- Use the 'recruit' command: recruit <type> <count>
- Each villager costs food to maintain
- You need sufficient housing (huts) to support your population
- Housing capacity is displayed in the Status panel

[#f1c40f]Assigning Villagers:[#ffffff]
- Use the 'assign' command: assign <type> <resource> <count>
- Available resources depend on your current age
- Unassigned villagers are considered "idle"
- You can unassign villagers with: unassign <type> <resource> <count>

[#f1c40f]Villager Types:[#ffffff]
- Basic villager: Available from the start
- Scholar: Unlocked in the Medieval Age, generates more knowledge

[#f1c40f]Tips:[#ffffff]
- Balance your villager assignments based on your current needs
- Make sure to maintain food production to support your population
- Build enough huts early to allow population growth
- In early ages, focus on resources needed for advancement
`)

	ls.AddTopic("resources", "Resources Guide", `
[::b][#3498db]RESOURCES[#ffffff][::b]

Resources are the foundation of your civilization's growth and advancement.
Different resources become available as you progress through ages.

[#f1c40f]Basic Resources:[#ffffff]
- Food: Required to support your population and recruit villagers
- Wood: Used for basic buildings and early technology
- Stone: Used for more advanced buildings (Bronze Age+)
- Gold: Used for trade and advanced structures (Iron Age+)
- Knowledge: Used for research and technological advancement (Iron Age+)

[#f1c40f]Resource Gathering:[#ffffff]
- Assign villagers to gather resources using: assign villager <resource> <count>
- Buildings can provide passive resource generation
- Resources have different gathering rates based on age and technology

[#f1c40f]Resource Management Tips:[#ffffff]
- Food should be your initial priority to support population growth
- Balance wood gathering for building construction
- Once in Bronze Age, focus on stone for advancing to Iron Age
- In Iron Age, balance gold and knowledge for further advancement
- Check age advancement requirements (status command) to plan resource gathering
`)

	ls.AddTopic("buildings", "Buildings Guide", `
[::b][#3498db]BUILDINGS[#ffffff][::b]

Buildings provide various benefits to your civilization, from housing to
resource production to research capabilities.

[#f1c40f]Available Buildings:[#ffffff]
- Hut: Provides housing for 2 villagers (Stone Age)
- Farm: Produces food passively (Stone Age)
- Lumber Mill: Increases wood production (Bronze Age)
- Mine: Produces stone and some gold (Bronze Age)
- Market: Generates gold (Iron Age)
- Library: Generates knowledge (Iron Age)

[#f1c40f]Building Construction:[#ffffff]
- Use the 'build' command: build <building>
- Each building has a resource cost (see 'buildings' command)
- Buildings provide passive benefits each tick
- Some buildings unlock new capabilities or resources

[#f1c40f]Building Strategy:[#ffffff]
- Build huts early to increase population capacity
- Farms provide passive food, freeing villagers for other tasks
- Balance building types based on your current resource needs
- Buildings are required for age advancement
- Use the 'buildings' command to see costs and available buildings
`)

	ls.AddTopic("ages", "Age Progression Guide", `
[::b][#3498db]AGE PROGRESSION[#ffffff][::b]

Your civilization advances through multiple ages, each unlocking new
buildings, resources, and capabilities.

[#f1c40f]Age Progression:[#ffffff]
- Stone Age: Starting age, basic resources and buildings
- Bronze Age: Unlocks stone, lumber mills, and mines
- Iron Age: Unlocks gold, knowledge, markets, and libraries
- Medieval Age: Unlocks scholars and advanced technologies
- Renaissance Age: Unlocks advanced economic capabilities
- Industrial Age: Unlocks mass production capabilities
- Modern Age: Final age with the most advanced technologies

[#f1c40f]Advancement Requirements:[#ffffff]
Each age has specific resource and building requirements to advance.
Use the 'status' command to see requirements for the next age.

[#f1c40f]Advancement Strategy:[#ffffff]
- Focus on meeting the specific requirements for the next age
- Balance immediate needs with long-term advancement goals
- Each new age brings significant advantages, prioritize advancement
- New buildings and resources in each age create new opportunities
`)

	ls.AddTopic("commands", "Game Commands", `
[::b][#3498db]GAME COMMANDS[#ffffff][::b]

CivIdleCli uses text commands to control all aspects of gameplay.

[#f1c40f]Basic Commands:[#ffffff]
- help: Display all available commands
- status: Show detailed status of your civilization
- quit: Exit the game

[#f1c40f]Resource Commands:[#ffffff]
- gather <resource> <count>: Quick assign villagers to gather
- assign <villager_type> <resource> <count>: Assign villagers to tasks
- unassign <villager_type> <resource> <count>: Unassign villagers

[#f1c40f]Building Commands:[#ffffff]
- build <building>: Build a structure
- buildings: List available buildings and their costs

[#f1c40f]Villager Commands:[#ffffff]
- recruit <villager_type> <count>: Recruit new villagers

[#f1c40f]Research Commands:[#ffffff]
- research <technology>: Start researching a technology
- techs: List available technologies for research

[#f1c40f]Save/Load Commands:[#ffffff]
- save <filename>: Save the current game
- load <filename>: Load a saved game
- saves: List all saved games

[#f1c40f]Other Commands:[#ffffff]
- stats: Display game statistics
- clear: Clear the console screen
- library <topic>: Access the help library (you are using it now!)
`)

	ls.AddTopic("tips", "Gameplay Tips & Strategies", `
[::b][#3498db]GAMEPLAY TIPS & STRATEGIES[#ffffff][::b]

Here are some helpful tips to succeed in CivIdleCli:

[#f1c40f]Early Game (Stone Age):[#ffffff]
- Focus on building 4-5 huts first to increase population capacity
- Recruit villagers up to your capacity
- Assign most villagers to gathering food and wood
- Build a farm to generate passive food
- Work toward Bronze Age requirements

[#f1c40f]Mid Game (Bronze/Iron Age):[#ffffff]
- Balance resource gathering based on advancement needs
- Build specialized production buildings (mines, lumber mills)
- Start focusing on knowledge once libraries are available
- Don't neglect food production as your population grows
- Gradually shift toward gold production for later ages

[#f1c40f]Late Game (Medieval+):[#ffffff]
- Utilize scholars to accelerate knowledge production
- Focus on a balanced economy with all resource types
- Research technologies that enhance production
- Maintain a reserve of resources for unexpected needs
- Push toward the final ages for maximum capabilities

[#f1c40f]General Tips:[#ffffff]
- Use the 'status' command frequently to check advancement requirements
- Keep some villagers free for quick reassignment as needs change
- Build multiple of each production building for compounding benefits
- Don't leave villagers idle - they consume food but produce nothing
- Save your game periodically using the 'save' command
`)

	return ls
}

// AddTopic adds a new topic to the library
func (ls *LibrarySystem) AddTopic(id, title, content string) {
	ls.topics[id] = &LibraryTopic{
		Title:   title,
		Content: content,
	}
}

// GetTopic returns a specific topic by ID
func (ls *LibrarySystem) GetTopic(id string) *LibraryTopic {
	return ls.topics[id]
}

// GetTopicList returns a list of all available topics
func (ls *LibrarySystem) GetTopicList() map[string]string {
	result := make(map[string]string)
	for id, topic := range ls.topics {
		result[id] = topic.Title
	}
	return result
}

// SearchTopics searches for topics containing the query in title or content
func (ls *LibrarySystem) SearchTopics(query string) map[string]string {
	result := make(map[string]string)
	query = strings.ToLower(query)

	for id, topic := range ls.topics {
		if strings.Contains(strings.ToLower(topic.Title), query) ||
			strings.Contains(strings.ToLower(topic.Content), query) {
			result[id] = topic.Title
		}
	}

	return result
}

// FormatTopicList returns a formatted string of all topics
func (ls *LibrarySystem) FormatTopicList() string {
	var result strings.Builder

	result.WriteString("[::b][#3498db]LIBRARY TOPICS[#ffffff][::b]\n\n")

	for id, topic := range ls.topics {
		result.WriteString(fmt.Sprintf("[#f1c40f]%s[#ffffff] - %s\n", id, topic.Title))
	}

	result.WriteString("\nUse 'library <topic>' to view a specific topic.")

	return result.String()
}
