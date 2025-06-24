package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/civcli/game"
)

// Dashboard provides the main game interface
type Dashboard struct {
	ui   *UIManager
	view *tview.Flex

	// Layout panels
	statsPanel     *tview.TextView
	buildingsPanel *tview.TextView
	researchPanel  *tview.TextView
	logPanel       *tview.TextView
	commandInput   *tview.InputField
	helpText       *tview.TextView

	// State tracking
	gameState *game.GameState
	messages  []Message
}

// Message represents a game message with type and timestamp
type Message struct {
	Text      string
	Type      string // "info", "success", "warning", "error"
	Timestamp time.Time
}

// NewDashboard creates a new game dashboard
func NewDashboard(ui *UIManager) *Dashboard {
	d := &Dashboard{
		ui:       ui,
		view:     tview.NewFlex(),
		messages: make([]Message, 0),
	}

	d.initializePanels()
	d.setupLayout()
	d.setupInputHandling()

	return d
}

// initializePanels creates all dashboard components
func (d *Dashboard) initializePanels() {
	theme := d.ui.GetTheme()

	// Stats panel - top left
	d.statsPanel = tview.NewTextView()
	d.statsPanel.SetBorder(true).
		SetTitle(" 📊 Civilization Stats ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	d.statsPanel.SetDynamicColors(true)

	// Buildings panel - top right
	d.buildingsPanel = tview.NewTextView()
	d.buildingsPanel.SetBorder(true).
		SetTitle(" 🏛️ Buildings & Infrastructure ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	d.buildingsPanel.SetDynamicColors(true)

	// Research panel - middle right
	d.researchPanel = tview.NewTextView()
	d.researchPanel.SetBorder(true).
		SetTitle(" 🔬 Research & Technology ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	d.researchPanel.SetDynamicColors(true)

	// Log panel - bottom
	d.logPanel = tview.NewTextView()
	d.logPanel.SetBorder(true).
		SetTitle(" 📜 Recent Events ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	d.logPanel.SetDynamicColors(true).
		SetScrollable(true)

	// Command input - bottom
	d.commandInput = tview.NewInputField().
		SetLabel("Command: ").
		SetPlaceholder("Type commands here (e.g., 'build house', 'research farming') or 'help' for assistance").
		SetFieldBackgroundColor(theme.Background).
		SetFieldTextColor(theme.Foreground)

	// Help text - bottom
	d.helpText = tview.NewTextView().
		SetText(" Press [yellow]F1[white] for help • [yellow]Ctrl+Q[white] to quit • [yellow]Tab[white] to navigate ").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)

	// Initialize with default content
	d.updateStatsDisplay()
	d.updateBuildingsDisplay()
	d.updateResearchDisplay()
	d.updateLogDisplay()
}

// setupLayout arranges the dashboard components
func (d *Dashboard) setupLayout() {
	// Create main horizontal split
	mainSplit := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Left column (stats and log)
	leftColumn := tview.NewFlex().SetDirection(tview.FlexRow)
	leftColumn.
		AddItem(d.statsPanel, 0, 1, false).
		AddItem(d.logPanel, 0, 1, false)

	// Right column (buildings and research)
	rightColumn := tview.NewFlex().SetDirection(tview.FlexRow)
	rightColumn.
		AddItem(d.buildingsPanel, 0, 1, false).
		AddItem(d.researchPanel, 0, 1, false)

	// Add columns to main split
	mainSplit.
		AddItem(leftColumn, 0, 1, false).
		AddItem(rightColumn, 0, 1, false)

	// Bottom input area
	inputArea := tview.NewFlex().SetDirection(tview.FlexRow)
	inputArea.
		AddItem(d.commandInput, 1, 0, true).
		AddItem(d.helpText, 1, 0, false)

	// Main layout
	d.view.SetDirection(tview.FlexRow)
	d.view.
		AddItem(mainSplit, 0, 1, false).
		AddItem(inputArea, 3, 0, true)
}

// setupInputHandling configures command input processing
func (d *Dashboard) setupInputHandling() {
	d.commandInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := d.commandInput.GetText()
			if command != "" {
				// Send command to game engine
				d.ui.SendInput(command)
				d.addMessage(fmt.Sprintf("> %s", command), "command")
				d.commandInput.SetText("")
			}
		}
	})

	// Set up input capture for navigation
	d.commandInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			// Tab cycles through focusable elements
			return event
		}
		return event
	})
}

// UpdateState updates the dashboard with new game state
func (d *Dashboard) UpdateState(state game.GameState) {
	d.gameState = &state
	d.updateStatsDisplay()
	d.updateBuildingsDisplay()
	d.updateResearchDisplay()
	// Remove the direct Draw() call to prevent potential deadlocks
}

// updateStatsDisplay refreshes the stats panel
func (d *Dashboard) updateStatsDisplay() {
	var content strings.Builder

	if d.gameState != nil {
		state := d.gameState

		content.WriteString(fmt.Sprintf("[yellow]📅 Age:[white] %s\n", state.Age))
		content.WriteString(fmt.Sprintf("[yellow]⏰ Tick:[white] %d\n", state.Tick))
		content.WriteString(fmt.Sprintf("[yellow]👥 Villagers:[white] %d/%d\n", len(state.Villagers), state.VillagerCap))
		content.WriteString("\n[cyan]Resources:[white]\n")

		// Display resources from the map
		for resource, amount := range state.Resources {
			emoji := d.getResourceEmoji(resource)
			content.WriteString(fmt.Sprintf("  %s %s: %.1f\n", emoji, resource, amount))
		}

		content.WriteString(fmt.Sprintf("\n[green]Total Food:[white] %.1f\n", state.TotalFood))
		content.WriteString(fmt.Sprintf("[yellow]Tick Duration:[white] %.1fs\n", state.TickDurationSeconds))
	} else {
		content.WriteString("[yellow]Welcome to CivIdleCli![white]\n\n")
		content.WriteString("🌟 Starting a new civilization...\n")
		content.WriteString("🏕️ Stone Age begins\n")
		content.WriteString("👥 Starting villagers...\n")
		content.WriteString("🌾 Gathering food...\n")
		content.WriteString("\n[cyan]Tip:[white] Start by building")
		content.WriteString(" houses to support more villagers!")
	}

	d.statsPanel.SetText(content.String())
}

// updateBuildingsDisplay refreshes the buildings panel
func (d *Dashboard) updateBuildingsDisplay() {
	var content strings.Builder

	if d.gameState != nil && len(d.gameState.Buildings) > 0 {
		content.WriteString("[cyan]Current Buildings:[white]\n\n")
		for building, count := range d.gameState.Buildings {
			content.WriteString(fmt.Sprintf("🏠 %s: %d\n", building, count))
		}
	} else {
		content.WriteString("[yellow]Available Buildings:[white]\n\n")
		content.WriteString("🏠 [green]House[white] - Increase population capacity\n")
		content.WriteString("   Cost: 10 Wood, 5 Stone\n\n")
		content.WriteString("🌾 [green]Farm[white] - Produce food automatically\n")
		content.WriteString("   Cost: 5 Wood, 15 Food\n\n")
		content.WriteString("🌲 [green]Lumber Mill[white] - Produce wood\n")
		content.WriteString("   Cost: 20 Wood, 10 Stone\n\n")
		content.WriteString("⛏️ [green]Quarry[white] - Produce stone\n")
		content.WriteString("   Cost: 15 Wood, 25 Stone\n\n")
		content.WriteString("[cyan]Use:[white] 'build <building>' to construct")
	}

	d.buildingsPanel.SetText(content.String())
}

// updateResearchDisplay refreshes the research panel
func (d *Dashboard) updateResearchDisplay() {
	var content strings.Builder

	if d.gameState != nil {
		if len(d.gameState.Research.Researched) > 0 {
			content.WriteString("[cyan]Researched Technologies:[white]\n\n")
			for _, tech := range d.gameState.Research.Researched {
				content.WriteString(fmt.Sprintf("🔬 %s\n", tech))
			}
			content.WriteString("\n")
		}

		if d.gameState.Research.Current != "" {
			content.WriteString(fmt.Sprintf("[yellow]Current Research:[white] %s\n", d.gameState.Research.Current))
			progress := (d.gameState.Research.Progress / d.gameState.Research.Cost) * 100
			content.WriteString(fmt.Sprintf("[cyan]Progress:[white] %.1f%%\n\n", progress))
		}

		content.WriteString("[yellow]Available Research:[white]\n")
		content.WriteString("🔬 [green]Agriculture[white] - Unlock advanced farming\n")
		content.WriteString("⚒️ [green]Tool Making[white] - Create better tools\n")
		content.WriteString("🏛️ [green]Construction[white] - Build larger structures\n")
	} else {
		content.WriteString("[yellow]Available Research:[white]\n\n")
		content.WriteString("🔬 [green]Agriculture[white] - Unlock advanced farming\n")
		content.WriteString("   Improves food production\n\n")
		content.WriteString("⚒️ [green]Tool Making[white] - Create better tools\n")
		content.WriteString("   Improves resource gathering\n\n")
		content.WriteString("🏛️ [green]Construction[white] - Build larger structures\n")
		content.WriteString("   Unlock new building types\n\n")
	}

	content.WriteString("\n[cyan]Use:[white] 'research <technology>' to advance")
	d.researchPanel.SetText(content.String())
}

// updateLogDisplay refreshes the message log
func (d *Dashboard) updateLogDisplay() {
	var content strings.Builder

	if len(d.messages) == 0 {
		content.WriteString("[cyan]Game started successfully![white]\n")
		content.WriteString("Type 'help' for a list of available commands.\n")
		content.WriteString("Begin by building houses and farms to grow your civilization!\n")
	} else {
		// Show last 10 messages
		start := 0
		if len(d.messages) > 10 {
			start = len(d.messages) - 10
		}

		for i := start; i < len(d.messages); i++ {
			msg := d.messages[i]
			color := d.getMessageColor(msg.Type)
			timestamp := msg.Timestamp.Format("15:04")
			content.WriteString(fmt.Sprintf("[gray]%s[white] %s%s[white]\n",
				timestamp, color, msg.Text))
		}
	}

	d.logPanel.SetText(content.String())
}

// getResourceEmoji returns an appropriate emoji for the resource type
func (d *Dashboard) getResourceEmoji(resource string) string {
	switch resource {
	case "food", "foraging", "hunting":
		return "🌾"
	case "wood", "lumber":
		return "🪵"
	case "stone":
		return "🗿"
	case "tools":
		return "⚒️"
	case "population":
		return "👥"
	default:
		return "📦"
	}
}

// getMessageColor returns the appropriate color for message type
func (d *Dashboard) getMessageColor(msgType string) string {
	switch msgType {
	case "success":
		return "[green]"
	case "warning":
		return "[yellow]"
	case "error":
		return "[red]"
	case "command":
		return "[cyan]"
	default:
		return "[white]"
	}
}

// ShowMessage adds a message to the log
func (d *Dashboard) ShowMessage(message, msgType string) {
	d.addMessage(message, msgType)
	d.updateLogDisplay()
	// Remove Draw() call to prevent potential conflicts during page transitions
}

// addMessage adds a message to the internal log
func (d *Dashboard) addMessage(message, msgType string) {
	msg := Message{
		Text:      message,
		Type:      msgType,
		Timestamp: time.Now(),
	}

	d.messages = append(d.messages, msg)

	// Keep only last 50 messages
	if len(d.messages) > 50 {
		d.messages = d.messages[len(d.messages)-50:]
	}
}

// GetView returns the dashboard view
func (d *Dashboard) GetView() tview.Primitive {
	return d.view
}

// Focus sets focus to the dashboard
func (d *Dashboard) Focus() {
	d.ui.GetApp().SetFocus(d.commandInput)
}
