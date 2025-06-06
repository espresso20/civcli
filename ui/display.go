package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/civcli/game"
	"github.com/user/civcli/ui/components"
)

// Corrected ShowSplashScreen function
func ShowSplashScreen() error {
	return components.ShowSplashScreen()
}

type Display struct {
	app          *tview.Application
	pages        *tview.Pages
	dashboard    *tview.Flex
	Resources    *tview.TextView
	Villagers    *tview.TextView
	Buildings    *tview.TextView
	Status       *tview.TextView
	Research     *tview.TextView
	Output       *tview.TextView
	CommandInput *tview.InputField
	mu           sync.Mutex
}

func NewDisplay() *Display {
	d := &Display{
		app:          tview.NewApplication(),
		pages:        tview.NewPages(),
		Resources:    tview.NewTextView(),
		Villagers:    tview.NewTextView(),
		Buildings:    tview.NewTextView(),
		Status:       tview.NewTextView(),
		Research:     tview.NewTextView(),
		Output:       tview.NewTextView(),
		CommandInput: tview.NewInputField().SetLabel("> ").SetFieldWidth(0),
	}

	d.Resources.SetBorder(true).SetTitle("Resources")
	d.Villagers.SetBorder(true).SetTitle("Villagers")
	d.Buildings.SetBorder(true).SetTitle("Buildings")
	d.Status.SetBorder(true).SetTitle("Status")
	d.Research.SetBorder(true).SetTitle("Research")
	d.Output.SetBorder(true).SetTitle("Output")
	d.Output.SetScrollable(true)

	// Layout: left (resources, villagers, buildings), right (status, research), bottom (output, input)
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
	d.pages.AddPage("dashboard", d.dashboard, true, true)

	// Add global key bindings
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Ctrl+C to quit (in addition to Ctrl+Q)
		if event.Key() == tcell.KeyCtrlC {
			d.app.Stop()
			return nil
		}
		// Ctrl+Q to quit
		if event.Key() == tcell.KeyCtrlQ {
			d.app.Stop()
			return nil
		}
		// Tab to toggle focus between output and input
		if event.Key() == tcell.KeyTab {
			if d.app.GetFocus() == d.CommandInput {
				d.app.SetFocus(d.Output)
			} else {
				d.app.SetFocus(d.CommandInput)
			}
			return nil
		}
		return event
	})

	d.app.SetRoot(d.pages, true)

	return d
}

func (d *Display) Start() error {
	return d.app.Run()
}

// ShowHelp displays the help screen with available commands
func (d *Display) ShowHelp(commands map[string]string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.app.QueueUpdateDraw(func() {
		components.ShowHelpWithApp(d.app, d.pages, commands, "dashboard")
	})
}

func (d *Display) Stop() {
	d.app.Stop()
}

// Example: update panels (call these from your game loop)
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

// --- Game DisplayInterface stubs ---
// You can fill these in with real logic later.

// DisplayDashboard updates all panels with game state data
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
	if state.Research.Current != "" {
		researchText = fmt.Sprintf("Current: %s\nProgress: %.1f/%.1f",
			state.Research.Current, state.Research.Progress, state.Research.Cost)
	} else {
		researchText = "Current: None"
	}

	if len(state.Research.Researched) > 0 {
		researchText += "\n\nCompleted:"
		for _, r := range state.Research.Researched {
			researchText += "\n- " + r
		}
	} else {
		researchText += "\n\nCompleted: None"
	}
	d.SetResearch(researchText)
}

func (d *Display) ShowMessage(message string, style string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Format the message based on style
	formattedMessage := components.FormatMessage(message, style)

	// Set the output text with the formatted message
	d.app.QueueUpdateDraw(func() {
		components.AppendToOutput(d.Output, formattedMessage)
	})
}

func (d *Display) ShowAgeAdvancement(newAge string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.app.QueueUpdateDraw(func() {
		components.ShowAgeAdvancement(d.pages, newAge, "dashboard")
	})

	// Also add to the output log
	d.ShowMessage("Advanced to "+newAge, "success")
}

func (d *Display) ShowLibraryContent(title, content string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.app.QueueUpdateDraw(func() {
		components.ShowLibraryContent(d.pages, title, content, "dashboard")
	})
}

func (d *Display) ShowLibraryTopicsList(topics map[string]string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.app.QueueUpdateDraw(func() {
		components.ShowLibraryTopicsList(d.pages, topics, "dashboard")
	})
}

func (d *Display) GetInput() (string, error) {
	// Create a channel to receive the command
	ch := make(chan string, 1)

	// Set up the input field
	d.app.QueueUpdateDraw(func() {
		d.CommandInput.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				command := d.CommandInput.GetText()
				if command != "" {
					ch <- command
					d.CommandInput.SetText("")
				}
			}
		})

		// Make sure the command input has focus
		d.app.SetFocus(d.CommandInput)
	})

	// Return the command from the channel
	return <-ch, nil
}
