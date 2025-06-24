package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/civcli/game"
)

// Simple, clean Display struct with only essential components
type Display struct {
	app          *tview.Application
	dashboard    *tview.Flex
	Resources    *tview.TextView
	Villagers    *tview.TextView
	Buildings    *tview.TextView
	Status       *tview.TextView
	Research     *tview.TextView
	Output       *tview.TextView
	CommandInput *tview.InputField
	inputChan    chan string
	mu           sync.Mutex
}

// NewDisplay creates a new, simple display with just the dashboard
func NewDisplay() *Display {
	d := &Display{
		app:          tview.NewApplication(),
		Resources:    tview.NewTextView(),
		Villagers:    tview.NewTextView(),
		Buildings:    tview.NewTextView(),
		Status:       tview.NewTextView(),
		Research:     tview.NewTextView(),
		Output:       tview.NewTextView(),
		CommandInput: tview.NewInputField().SetLabel("> ").SetFieldWidth(0),
		inputChan:    make(chan string, 10),
	}

	// Configure panels
	d.Resources.SetBorder(true).SetTitle("Resources")
	d.Villagers.SetBorder(true).SetTitle("Villagers")
	d.Buildings.SetBorder(true).SetTitle("Buildings")
	d.Status.SetBorder(true).SetTitle("Status")
	d.Research.SetBorder(true).SetTitle("Research")
	d.Output.SetBorder(true).SetTitle("Output")
	d.Output.SetScrollable(true)

	// Create simple layout
	leftCol := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.Resources, 0, 1, false).
		AddItem(d.Villagers, 0, 1, false).
		AddItem(d.Buildings, 0, 1, false)
	rightCol := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(d.Status, 0, 1, false).
		AddItem(d.Research, 0, 1, false)
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(leftCol, 0, 2, false).
		AddItem(rightCol, 0, 3, false)
	d.dashboard = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(topRow, 0, 3, false).
		AddItem(d.Output, 0, 1, false).
		AddItem(d.CommandInput, 1, 0, true)

	// Set up simple input handling
	d.setupInputHandler()

	// Simple global key bindings - just quit keys
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ {
			d.app.Stop()
			return nil
		}
		return event
	})

	d.app.SetRoot(d.dashboard, true)
	return d
}

// setupInputHandler sets up simple, reliable input handling
func (d *Display) setupInputHandler() {
	d.CommandInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			command := d.CommandInput.GetText()
			if command != "" {
				select {
				case d.inputChan <- command:
					d.CommandInput.SetText("")
				default:
					// Channel full, ignore
				}
			}
		}
	})
}

// Start runs the application
func (d *Display) Start() error {
	return d.app.Run()
}

// Stop stops the application
func (d *Display) Stop() {
	d.app.Stop()
}

// GetInput gets user input with timeout (non-blocking)
func (d *Display) GetInput() (string, error) {
	d.app.QueueUpdateDraw(func() {
		d.app.SetFocus(d.CommandInput)
	})

	select {
	case command := <-d.inputChan:
		return command, nil
	case <-time.After(100 * time.Millisecond):
		return "", fmt.Errorf("no input ready")
	}
}

// Simple panel update methods
func (d *Display) SetResources(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Resources.SetText(text) })
}

func (d *Display) SetVillagers(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Villagers.SetText(text) })
}

func (d *Display) SetBuildings(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Buildings.SetText(text) })
}

func (d *Display) SetStatus(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Status.SetText(text) })
}

func (d *Display) SetResearch(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Research.SetText(text) })
}

func (d *Display) SetOutput(text string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.app.QueueUpdateDraw(func() { d.Output.SetText(text) })
}

// ShowMessage displays a simple message in the output
func (d *Display) ShowMessage(message, msgType string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var coloredMessage string
	switch msgType {
	case "error":
		coloredMessage = fmt.Sprintf("[red]%s[white]", message)
	case "success":
		coloredMessage = fmt.Sprintf("[green]%s[white]", message)
	case "warning":
		coloredMessage = fmt.Sprintf("[yellow]%s[white]", message)
	default:
		coloredMessage = message
	}

	d.app.QueueUpdateDraw(func() {
		current := d.Output.GetText(false)
		d.Output.SetText(current + coloredMessage + "\n")
	})
}

// DisplayDashboard updates all panels with game state
func (d *Display) DisplayDashboard(state game.GameState) {
	// Format resources
	resourcesText := ""
	if len(state.Resources) > 0 {
		for name, amount := range state.Resources {
			resourcesText += fmt.Sprintf("%s: %.1f\n", name, amount)
		}
	} else {
		resourcesText = "No resources available"
	}
	d.SetResources(resourcesText)

	// Format villagers
	villagersText := ""
	totalVillagers := 0
	if len(state.Villagers) > 0 {
		for name, info := range state.Villagers {
			villagersText += fmt.Sprintf("%s: %d\n", name, info.Count)
			totalVillagers += info.Count
		}
	} else {
		villagersText = "No villagers available"
	}
	d.SetVillagers(villagersText)

	// Format buildings
	buildingsText := ""
	if len(state.Buildings) > 0 {
		for name, count := range state.Buildings {
			buildingsText += fmt.Sprintf("%s: %d\n", name, count)
		}
	} else {
		buildingsText = "No buildings available"
	}
	d.SetBuildings(buildingsText)

	// Format status
	statusText := fmt.Sprintf("Age: %s\nTick: %d\nVillagers: %d/%d\nTick Duration: %.1f sec",
		state.Age, state.Tick, totalVillagers, state.VillagerCap, state.TickDurationSeconds)
	d.SetStatus(statusText)

	// Format research
	researchText := ""
	if len(state.Research.Researched) > 0 {
		researchText = "Completed Technologies:\n"
		for _, tech := range state.Research.Researched {
			researchText += fmt.Sprintf("• %s\n", tech)
		}
	} else {
		researchText = "No technologies researched"
	}

	if state.Research.Current != "" {
		researchText += fmt.Sprintf("\nCurrently researching: %s (%.1f/%.1f)",
			state.Research.Current, state.Research.Progress, state.Research.Cost)
	}

	d.SetResearch(researchText)
}

// ShowHelp displays a simple help message (no complex system)
func (d *Display) ShowHelp(commands map[string]string) {
	helpText := "Available Commands:\n"
	helpText += "==================\n\n"

	for cmd, desc := range commands {
		helpText += fmt.Sprintf("%s: %s\n", cmd, desc)
	}

	helpText += "\nPress Ctrl+C or Ctrl+Q to quit."

	d.ShowMessage(helpText, "info")
}

// ShowAgeAdvancement shows advancement notification
func (d *Display) ShowAgeAdvancement(newAge string) {
	message := fmt.Sprintf("Civilization advanced to %s!", newAge)
	d.ShowMessage(message, "success")
}

// NotifyAdvancement shows advancement notification (legacy method)
func (d *Display) NotifyAdvancement(oldAge, newAge string) {
	message := fmt.Sprintf("Civilization advanced from %s to %s!", oldAge, newAge)
	d.ShowMessage(message, "success")
}
