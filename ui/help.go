package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HelpSystem provides comprehensive in-game help and tutorials
type HelpSystem struct {
	ui         *UIManager
	view       *tview.Flex
	navigation *tview.List
	content    *tview.TextView
	returnPage string
}

// NewHelpSystem creates a new help system
func NewHelpSystem(ui *UIManager) *HelpSystem {
	h := &HelpSystem{
		ui:         ui,
		view:       tview.NewFlex(),
		navigation: tview.NewList(),
		content:    tview.NewTextView(),
		returnPage: "splash",
	}

	h.setupNavigation()
	h.setupContent()
	h.setupLayout()

	return h
}

// setupNavigation creates the help navigation menu
func (h *HelpSystem) setupNavigation() {
	theme := h.ui.GetTheme()

	h.navigation.
		SetBorder(true).
		SetTitle(" 📚 Help Topics ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)

	// Add help topics
	h.navigation.AddItem("🚀 Getting Started", "Learn the basics of the game", '1', func() {
		h.showGettingStarted()
	})

	h.navigation.AddItem("💼 Commands", "List of all available commands", '2', func() {
		h.showCommands()
	})

	h.navigation.AddItem("🏛️ Buildings", "Information about buildings", '3', func() {
		h.showBuildings()
	})

	h.navigation.AddItem("🔬 Research", "Technology and research system", '4', func() {
		h.showResearch()
	})

	h.navigation.AddItem("📊 Resources", "Resource management guide", '5', func() {
		h.showResources()
	})

	h.navigation.AddItem("👥 Villagers", "Population and workforce management", '6', func() {
		h.showVillagers()
	})

	h.navigation.AddItem("⏰ Game Mechanics", "Understanding ticks and progression", '7', func() {
		h.showGameMechanics()
	})

	h.navigation.AddItem("💡 Tips & Strategy", "Advanced gameplay tips", '8', func() {
		h.showTipsAndStrategy()
	})

	h.navigation.AddItem("🔙 Return to Game", "Go back to the previous screen", 'q', func() {
		h.ui.HideHelp()
	})

	// Set up navigation input capture
	h.navigation.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			h.ui.HideHelp()
			return nil
		}
		return event
	})
}

// setupContent creates the help content area
func (h *HelpSystem) setupContent() {
	theme := h.ui.GetTheme()

	h.content.SetBorder(true).
		SetTitle(" 📖 Help Content ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)

	h.content.SetDynamicColors(true).
		SetWordWrap(true).
		SetScrollable(true)

	// Set up content input capture for scrolling
	h.content.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			h.ui.HideHelp()
			return nil
		}
		return event
	})

	// Show default content
	h.showGettingStarted()
}

// setupLayout arranges the help system components
func (h *HelpSystem) setupLayout() {
	h.view.SetDirection(tview.FlexColumn)

	h.view.
		AddItem(h.navigation, 30, 0, true). // Navigation sidebar
		AddItem(h.content, 0, 1, false)     // Main content area
}

// showGettingStarted displays the getting started guide
func (h *HelpSystem) showGettingStarted() {
	content := `[yellow::b]🚀 Getting Started with CivIdleCli[white::-]

Welcome to CivIdleCli, a command-line civilization building game! Here's everything you need to know to get started:

[cyan::b]📋 Basic Concepts[white::-]

• [yellow]Civilization:[white] You start in the Stone Age with a small group of villagers
• [yellow]Resources:[white] Gather food, wood, stone, and other materials to survive and grow
• [yellow]Buildings:[white] Construct structures to increase capacity and production
• [yellow]Research:[white] Advance through ages by researching new technologies
• [yellow]Villagers:[white] Your population - they gather resources and work in buildings

[cyan::b]🎯 Your First Steps[white::-]

1. [green]Check your status:[white] Look at the stats panel to see your current resources
2. [green]Build houses:[white] Type 'build house' to increase your population capacity
3. [green]Gather food:[white] Your villagers automatically forage and hunt for food
4. [green]Research agriculture:[white] Type 'research agriculture' to improve food production
5. [green]Expand:[white] Build more structures and research new technologies

[cyan::b]⌨️  Basic Commands[white::-]

• [yellow]build <building>[white] - Construct a building
• [yellow]research <technology>[white] - Research a new technology
• [yellow]status[white] - Show detailed civilization status
• [yellow]help[white] - Show this help system
• [yellow]save[white] - Save your current game
• [yellow]load[white] - Load a saved game

[green::b]💡 Pro Tip:[white::-] The game progresses automatically through "ticks" - you don't need to do anything for time to pass and resources to be gathered!

Press [yellow]ESC[white] or select another topic to continue learning.`

	h.content.SetText(content)
}

// showCommands displays all available commands
func (h *HelpSystem) showCommands() {
	content := `[yellow::b]💼 Complete Command Reference[white::-]

Here are all the commands available in CivIdleCli:

[cyan::b]🏗️ Building Commands[white::-]

• [green]build house[white] - Increase population capacity
• [green]build farm[white] - Automatic food production
• [green]build lumber mill[white] - Automatic wood production
• [green]build quarry[white] - Automatic stone production
• [green]build workshop[white] - Tool and equipment production

[cyan::b]🔬 Research Commands[white::-]

• [green]research agriculture[white] - Improve farming efficiency
• [green]research tool making[white] - Create better tools
• [green]research construction[white] - Unlock advanced buildings
• [green]research pottery[white] - Food storage improvements
• [green]research animal husbandry[white] - Livestock management

[cyan::b]📊 Information Commands[white::-]

• [green]status[white] - Show detailed civilization information
• [green]buildings[white] - List all buildings and their status
• [green]research[white] - Show research progress and available technologies
• [green]villagers[white] - Display villager information and assignments

[cyan::b]💾 Game Management[white::-]

• [green]save[white] - Save your current game progress
• [green]load[white] - Load a previously saved game
• [green]help[white] - Open this help system
• [green]quit[white] - Exit the game

[cyan::b]🎮 Shortcuts[white::-]

• [yellow]F1[white] - Quick help
• [yellow]Ctrl+Q[white] - Quick quit
• [yellow]Tab[white] - Navigate interface elements
• [yellow]ESC[white] - Return to previous screen

[green::b]💡 Command Tips:[white::-]
- Commands are case-insensitive
- You can use partial names (e.g., 'build h' for 'build house')
- Type 'help <command>' for detailed information about specific commands`

	h.content.SetText(content)
}

// showBuildings displays building information
func (h *HelpSystem) showBuildings() {
	content := `[yellow::b]🏛️ Buildings & Infrastructure[white::-]

Buildings are the foundation of your civilization. Each building type serves a specific purpose and requires certain resources to construct.

[cyan::b]🏠 Residential Buildings[white::-]

[green]House[white]
• Purpose: Increases population capacity
• Cost: 10 Wood, 5 Stone
• Capacity: +5 villagers per house
• Unlocked: Available from start

[cyan::b]🏭 Production Buildings[white::-]

[green]Farm[white]
• Purpose: Automatic food production
• Cost: 5 Wood, 15 Food (seeds)
• Production: +3 Food per tick
• Unlocked: Available from start

[green]Lumber Mill[white]
• Purpose: Automatic wood production
• Cost: 20 Wood, 10 Stone
• Production: +2 Wood per tick
• Unlocked: Available from start

[green]Quarry[white]
• Purpose: Automatic stone production
• Cost: 15 Wood, 25 Stone
• Production: +1 Stone per tick
• Unlocked: Available from start

[green]Workshop[white]
• Purpose: Tool and equipment production
• Cost: 30 Wood, 20 Stone, 10 Tools
• Production: +1 Tools per tick
• Unlocked: Requires Tool Making research

[cyan::b]🔬 Advanced Buildings[white::-]

[green]Granary[white]
• Purpose: Food storage and preservation
• Cost: 25 Wood, 15 Stone
• Effect: Reduces food spoilage by 50%
• Unlocked: Requires Pottery research

[green]Marketplace[white]
• Purpose: Resource trading and management
• Cost: 40 Wood, 30 Stone, 5 Tools
• Effect: Enables resource trading
• Unlocked: Requires Construction research

[cyan::b]🎯 Building Strategy[white::-]

1. [yellow]Start with houses[white] - More villagers = more resource gathering
2. [yellow]Build farms early[white] - Ensure food security for your growing population
3. [yellow]Balance production[white] - Don't focus on just one resource type
4. [yellow]Plan ahead[white] - Some buildings require resources from other buildings

[green::b]💡 Building Tips:[white::-]
• Buildings work automatically once constructed
• More buildings = faster resource generation
• Some buildings become more efficient with research upgrades`

	h.content.SetText(content)
}

// showResearch displays research information
func (h *HelpSystem) showResearch() {
	content := `[yellow::b]🔬 Research & Technology[white::-]

Research is how your civilization advances through the ages. Each technology unlocks new capabilities, buildings, or improves existing systems.

[cyan::b]🌾 Food & Agriculture[white::-]

[green]Agriculture[white]
• Effect: +50% farm efficiency, unlocks advanced crops
• Cost: 100 Research points
• Prerequisites: None
• Unlocks: Granary, advanced farming techniques

[green]Pottery[white]
• Effect: Food storage, reduced spoilage
• Cost: 150 Research points
• Prerequisites: Agriculture
• Unlocks: Granary, food preservation

[green]Animal Husbandry[white]
• Effect: Livestock for food and materials
• Cost: 200 Research points
• Prerequisites: Agriculture
• Unlocks: Livestock buildings, leather production

[cyan::b]⚒️ Tools & Crafting[white::-]

[green]Tool Making[white]
• Effect: +25% resource gathering efficiency
• Cost: 120 Research points
• Prerequisites: None
• Unlocks: Workshop, better tools

[green]Metalworking[white]
• Effect: Metal tools and weapons
• Cost: 300 Research points
• Prerequisites: Tool Making
• Unlocks: Bronze Age technologies

[cyan::b]🏗️ Construction & Engineering[white::-]

[green]Construction[white]
• Effect: Larger, more efficient buildings
• Cost: 180 Research points
• Prerequisites: None
• Unlocks: Marketplace, advanced buildings

[green]Architecture[white]
• Effect: Monumental buildings, city planning
• Cost: 400 Research points
• Prerequisites: Construction
• Unlocks: Temple, advanced city structures

[cyan::b]🎓 How Research Works[white::-]

1. [yellow]Generate Research Points:[white] Villagers automatically generate research
2. [yellow]Choose Technology:[white] Use 'research <technology>' command
3. [yellow]Wait for Completion:[white] Research progresses automatically over time
4. [yellow]Enjoy Benefits:[white] New capabilities unlock immediately

[cyan::b]📈 Research Strategy[white::-]

• [yellow]Early Game:[white] Focus on Agriculture and Tool Making
• [yellow]Mid Game:[white] Construction and specialized technologies
• [yellow]Late Game:[white] Advanced technologies for new ages

[green::b]💡 Research Tips:[white::-]
• Research continues even when you're not playing
• Some technologies are prerequisites for others
• Plan your research path based on your civilization's needs`

	h.content.SetText(content)
}

// showResources displays resource management information
func (h *HelpSystem) showResources() {
	content := `[yellow::b]📊 Resource Management Guide[white::-]

Resources are the lifeblood of your civilization. Understanding how to gather, manage, and spend them efficiently is key to success.

[cyan::b]🌾 Primary Resources[white::-]

[green]Food[white]
• Sources: Foraging, hunting, farms
• Usage: Population growth, building construction
• Storage: Unlimited (with spoilage)
• Tips: Build farms early, research agriculture

[green]Wood[white]
• Sources: Villager gathering, lumber mills
• Usage: All building construction
• Storage: Unlimited
• Tips: Most important early resource

[green]Stone[white]
• Sources: Villager gathering, quarries
• Usage: Advanced buildings, tools
• Storage: Unlimited
• Tips: Harder to gather, plan usage carefully

[green]Tools[white]
• Sources: Workshops, crafting
• Usage: Advanced buildings, research
• Storage: Unlimited
• Tips: Requires research to produce

[cyan::b]📈 Resource Generation[white::-]

[yellow]Automatic Gathering:[white]
• Villagers automatically gather resources each tick
• Gathering efficiency increases with better tools
• Population size affects total gathering rate

[yellow]Building Production:[white]
• Production buildings generate resources automatically
• Multiple buildings stack their effects
• Research can improve building efficiency

[cyan::b]⚖️ Resource Balance[white::-]

[green]Early Game (Stone Age):[white]
• Focus: Food security, basic shelter
• Priority: Food > Wood > Stone
• Strategy: Build houses and farms first

[green]Mid Game (Bronze Age):[white]
• Focus: Tool production, infrastructure
• Priority: Wood > Stone > Tools > Food
• Strategy: Diversify production, research key technologies

[green]Late Game (Iron Age+):[white]
• Focus: Advanced technologies, specialization
• Priority: Balanced production across all resources
• Strategy: Optimize efficiency, plan for next age

[cyan::b]💡 Resource Tips[white::-]

• [yellow]Monitor ratios:[white] Keep resources balanced for steady growth
• [yellow]Plan ahead:[white] Check building costs before starting construction
• [yellow]Invest in production:[white] Buildings pay for themselves over time
• [yellow]Research matters:[white] Technology improvements compound over time

[green::b]⚠️ Common Mistakes:[white::-]
• Overbuilding houses without food production
• Ignoring tool production for too long
• Not researching efficiency improvements
• Focusing on only one resource type`

	h.content.SetText(content)
}

// showVillagers displays villager management information
func (h *HelpSystem) showVillagers() {
	content := `[yellow::b]👥 Villager & Population Management[white::-]

Your villagers are the heart of your civilization. They gather resources, work in buildings, and research new technologies.

[cyan::b]🏠 Population Basics[white::-]

[green]Population Capacity:[white]
• Base capacity: 10 villagers
• Each house adds +5 capacity
• Growth stops when capacity is reached

[green]Population Growth:[white]
• Automatic growth over time with sufficient food
• Growth rate depends on available housing
• Food consumption increases with population

[cyan::b]⚒️ Villager Activities[white::-]

[yellow]Resource Gathering:[white]
• Food: Foraging, hunting (automatic)
• Wood: Tree cutting (automatic)
• Stone: Quarrying (automatic)
• Research: Knowledge generation (automatic)

[yellow]Building Operations:[white]
• Production buildings automatically assign workers
• More villagers = faster resource production
• Buildings work more efficiently with proper staffing

[cyan::b]📊 Efficiency Factors[white::-]

[green]Tools & Technology:[white]
• Better tools increase gathering efficiency
• Research unlocks new gathering methods
• Advanced buildings boost productivity

[green]Specialization:[white]
• Villagers become more efficient over time
• Balanced workforce prevents bottlenecks
• Technology can unlock specialist roles

[cyan::b]🎯 Population Strategy[white::-]

[yellow]Early Game:[white]
• Build houses to increase capacity
• Focus on food production to support growth
• Aim for 20-30 villagers quickly

[yellow]Mid Game:[white]
• Balance population with building construction
• Ensure adequate food production
• Consider specialist building assignments

[yellow]Late Game:[white]
• Optimize villager assignments for efficiency
• Use advanced buildings for maximum productivity
• Plan population for next age requirements

[cyan::b]💡 Population Tips[white::-]

• [green]Housing first:[white] Always build capacity before expecting growth
• [green]Food security:[white] More villagers need more food
• [green]Balanced growth:[white] Don't neglect infrastructure for population
• [green]Plan ahead:[white] Population growth takes time

[cyan::b]📈 Managing Large Populations[white::-]

[yellow]Efficiency Management:[white]
• Monitor resource consumption vs. production
• Use specialized buildings for better efficiency
• Research technologies that boost productivity

[yellow]Infrastructure Scaling:[white]
• Build more production buildings as population grows
• Ensure adequate resource storage
• Plan for increased complexity

[green::b]⚠️ Common Issues:[white::-]
• Building too many houses without food production
• Ignoring the relationship between population and consumption
• Not building enough production to support large populations`

	h.content.SetText(content)
}

// showGameMechanics displays game mechanics information
func (h *HelpSystem) showGameMechanics() {
	content := `[yellow::b]⏰ Game Mechanics & Tick System[white::-]

Understanding how CivIdleCli works under the hood will help you make better strategic decisions.

[cyan::b]🕐 The Tick System[white::-]

[green]What is a Tick?[white]
• A tick is one game cycle (default: 2 seconds)
• All automatic actions happen on each tick
• Resources are generated, research progresses, population grows

[green]Tick Duration:[white]
• Default: 2 seconds per tick
• Displayed in the stats panel
• Affects the pace of the game

[cyan::b]📈 Automatic Processes[white::-]

[yellow]Resource Generation:[white]
• Villagers gather resources automatically
• Buildings produce resources automatically
• Efficiency improvements apply immediately

[yellow]Research Progress:[white]
• Research points accumulate automatically
• Current research advances toward completion
• Multiple technologies can't be researched simultaneously

[yellow]Population Growth:[white]
• Population grows naturally over time
• Growth rate depends on available housing and food
• Growth stops when housing capacity is reached

[cyan::b]🎮 Real-Time vs. Idle[white::-]

[green]Active Play:[white]
• Issue commands to build and research
• Monitor progress and make strategic decisions
• React to resource shortages or opportunities

[green]Idle Progression:[white]
• The game continues even when you're not actively playing
• Resources accumulate automatically
• Research completes over time

[cyan::b]⚡ Action Types[white::-]

[yellow]Instant Actions:[white]
• Commands execute immediately when issued
• Building construction happens instantly (with resource cost)
• Research selection is immediate

[yellow]Progressive Actions:[white]
• Resource gathering happens over time
• Research completion takes multiple ticks
• Population growth occurs gradually

[cyan::b]🎯 Strategic Timing[white::-]

[green]Short-term Planning (1-10 ticks):[white]
• Monitor immediate resource needs
• Queue up quick construction projects
• Respond to resource shortages

[green]Medium-term Planning (10-50 ticks):[white]
• Plan research completion timing
• Coordinate building construction sequences
• Prepare for population growth spurts

[green]Long-term Planning (50+ ticks):[white]
• Plan age transitions
• Design civilization development path
• Optimize for maximum efficiency

[cyan::b]💡 Efficiency Tips[white::-]

• [yellow]Idle time is productive:[white] The game works for you
• [yellow]Plan during active time:[white] Make decisions, then let the game run
• [yellow]Check in regularly:[white] Monitor progress and adjust strategy
• [yellow]Compound growth:[white] Early investments pay off over many ticks

[green::b]🔄 Game Loop:[white::-]
1. Check current status and resources
2. Make strategic decisions (build, research)
3. Wait for automatic progress
4. Repeat with improved capabilities`

	h.content.SetText(content)
}

// showTipsAndStrategy displays advanced tips
func (h *HelpSystem) showTipsAndStrategy() {
	content := `[yellow::b]💡 Advanced Tips & Strategy[white::-]

Master these strategies to build the most efficient and successful civilization.

[cyan::b]🚀 Early Game Optimization[white::-]

[green]Opening Strategy:[white]
• Build 2-3 houses immediately (population capacity)
• Construct 1-2 farms for food security
• Start Agriculture research ASAP
• Monitor food consumption vs. production

[green]Resource Priority:[white]
1. [yellow]Food security[white] - Build farms before houses
2. [yellow]Population growth[white] - Houses enable more workers
3. [yellow]Wood production[white] - Lumber mills for construction materials
4. [yellow]Technology[white] - Research for efficiency improvements

[cyan::b]⚖️ Mid Game Balance[white::-]

[green]Production Scaling:[white]
• 1 farm per 10 villagers (approximate ratio)
• 1 lumber mill per 20 villagers
• 1 quarry per 30 villagers
• Adjust based on construction needs

[green]Research Path:[white]
• Agriculture → Tool Making → Construction
• Pottery for food storage improvements
• Metalworking to advance to Bronze Age

[cyan::b]📊 Advanced Optimization[white::-]

[yellow]Efficiency Multipliers:[white]
• Tool Making research: +25% gathering
• Agriculture research: +50% farm production
• Advanced buildings compound these effects

[yellow]Resource Stockpiling:[white]
• Keep 2-3 buildings worth of resources in reserve
• Plan for research costs (some are expensive)
• Save stone for advanced buildings

[yellow]Population Management:[white]
• Optimal population: 3-5x your building count
• Too many villagers without buildings = inefficiency
• Too few villagers = slow resource generation

[cyan::b]🎯 Age Transition Strategy[white::-]

[green]Stone to Bronze Age:[white]
• Requirements: Metalworking research, 50+ population
• Preparation: Stockpile 500+ of each basic resource
• Focus: Tool production, advanced buildings

[green]Bronze to Iron Age:[white]
• Requirements: Advanced technologies, established infrastructure
• Preparation: Diverse building types, efficient production
• Focus: Specialization, optimization

[cyan::b]💎 Expert Techniques[white::-]

[yellow]Resource Cycling:[white]
• Use excess food to fuel rapid building construction
• Convert basic resources into advanced ones through buildings
• Time research completion with resource availability

[yellow]Compound Growth:[white]
• Every building pays for itself over time
• Research improvements affect ALL future production
• Early investments have exponential returns

[yellow]Bottleneck Management:[white]
• Identify limiting resources (usually stone or tools)
• Build production specifically for bottleneck resources
• Research technologies that address constraints

[cyan::b]⚠️ Common Pitfalls[white::-]

• [red]Overbuilding houses[white] without corresponding food production
• [red]Ignoring research[white] - technology improvements are huge
• [red]Resource hoarding[white] - unused resources don't generate value
• [red]Impatience[white] - the game rewards planning and patience

[cyan::b]🏆 Mastery Goals[white::-]

[green]Efficiency Targets:[white]
• 100+ population by tick 200
• All basic research complete by Bronze Age
• Self-sustaining resource generation
• Positive growth in all resource types

[green]Advanced Challenges:[white]
• Reach Iron Age in minimal ticks
• Build the largest possible population
• Achieve maximum resource generation efficiency
• Complete all available research trees

[green::b]🎓 Remember:[white::-] The best strategy adapts to your civilization's current needs. Monitor, adjust, and optimize continuously!`

	h.content.SetText(content)
}

// GetView returns the help system view
func (h *HelpSystem) GetView() tview.Primitive {
	return h.view
}

// Focus sets focus to the help system
func (h *HelpSystem) Focus() {
	h.ui.GetApp().SetFocus(h.navigation)
}

// SetReturnPage sets which page to return to when help is closed
func (h *HelpSystem) SetReturnPage(page string) {
	h.returnPage = page
}

// GetReturnPage returns the page to return to when help is closed
func (h *HelpSystem) GetReturnPage() string {
	return h.returnPage
}
