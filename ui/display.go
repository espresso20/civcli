package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/user/civcli/game"
)

// Display handles all UI and display functionality
type Display struct {
	app        *tview.Application
	pages      *tview.Pages
	mainFlex   *tview.Flex
	dashboard  *tview.Grid
	input      *tview.InputField
	output     *tview.TextView
	resources  *tview.Table
	villagers  *tview.Table
	buildings  *tview.Table
	status     *tview.TextView
	research   *tview.TextView
	inputChan  chan string
	rightPanel *tview.Flex
}

// NewDisplay creates a new display
func NewDisplay() *Display {
	d := &Display{
		app:       tview.NewApplication(),
		pages:     tview.NewPages(),
		inputChan: make(chan string),
	}
	d.setupUI()
	return d
}

// setupUI sets up the UI components
func (d *Display) setupUI() {
	// Create main components
	d.output = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})
	d.output.SetBorder(true).SetTitle("Game Output")

	d.input = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0).
		SetFieldTextColor(tcell.ColorLightGray).   // Set input text color to light gray
		SetFieldBackgroundColor(tcell.ColorBlack). // Set input background to black
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := d.input.GetText()
				d.input.SetText("")
				// Only send to channel if it's not empty
				if text != "" {
					go func() {
						d.inputChan <- text
					}()
				}
			}
		})

	// Create dashboard components
	d.resources = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)
	d.resources.SetBorder(true).SetTitle("Resources")

	d.villagers = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)
	d.villagers.SetBorder(true).SetTitle("Villagers")

	d.buildings = tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)
	d.buildings.SetBorder(true).SetTitle("Buildings")

	d.status = tview.NewTextView().
		SetDynamicColors(true)
	d.status.SetBorder(true).SetTitle("Status")

	// Create research panel
	d.research = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			d.app.Draw()
		})
	d.research.SetBorder(true).SetTitle("Research")

	// Create dashboard layout with improved dimensions
	leftPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.resources, 10, 1, false).
		AddItem(d.buildings, 10, 1, false)

	d.rightPanel = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.villagers, 10, 1, false).
		AddItem(d.status, 5, 1, false).
		AddItem(d.research, 5, 1, false)

	topDashboard := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 1, false).
		AddItem(d.rightPanel, 0, 1, false)

	d.dashboard = tview.NewGrid().
		SetRows(20, 0).
		SetColumns(0).
		AddItem(topDashboard, 0, 0, 1, 1, 0, 0, false).
		AddItem(d.output, 1, 0, 1, 1, 0, 0, false)

	// Create main layout
	d.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.dashboard, 0, 1, false).
		AddItem(d.input, 1, 1, true)

	// Create intro page with improved styling
	introText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	introText.SetText(`
[::b][#2ecc71]CivIdleCli - A Command Line Civilization Builder Game[#ffffff][::b]

[#3498db]Welcome to CivIdleCli![#ffffff] You are the leader of a small tribe in the Stone Age.

[#f1c40f]Your Mission:[#ffffff]
- Gather resources to survive and grow
- Build structures to advance your civilization
- Research technologies to unlock new possibilities
- Assign villagers efficiently to different tasks
- Guide your people through the ages of history

[#e74c3c]Type 'help' at any time to see available commands.[#ffffff]

[::b][#2ecc71]Press Enter to begin your journey...[#ffffff][::b]
	`)

	introPage := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(introText, 15, 1, false).
			AddItem(nil, 0, 1, false), 70, 1, true).
		AddItem(nil, 0, 1, false)

	// Add pages
	d.pages.AddPage("intro", introPage, true, true)
	d.pages.AddPage("main", d.mainFlex, true, false)

	// Set up application
	d.app.SetRoot(d.pages, true).
		EnableMouse(true)

	// Set up key bindings for navigation
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Global hotkeys
		switch event.Key() {
		case tcell.KeyEsc:
			// Return focus to input
			d.app.SetFocus(d.input)
		case tcell.KeyF1:
			// Show help
			go func() {
				d.inputChan <- "help"
			}()
		}
		return event
	})
}

// ShowIntro shows the intro screen
func (d *Display) ShowIntro() {
	d.app.SetFocus(d.pages)

	// Set up a handler to transition from intro to main on Enter
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if d.pages.HasPage("intro") && event.Key() == tcell.KeyEnter {
			d.pages.SwitchToPage("main")
			// Keep global key capture but reset it
			d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				// Global hotkeys
				switch event.Key() {
				case tcell.KeyEsc:
					// Return focus to input
					d.app.SetFocus(d.input)
				case tcell.KeyF1:
					// Show help
					go func() {
						d.inputChan <- "help"
					}()
				}
				return event
			})
			d.app.SetFocus(d.input)
		}
		return event
	})

	// Start the app if not already running
	go func() {
		if err := d.app.Run(); err != nil {
			panic(err)
		}
	}()
}

// ShowHelp displays help information
func (d *Display) ShowHelp(commands map[string]string) {
	d.output.Clear()
	d.output.SetTextColor(tcell.ColorTeal)

	// Create a more organized help display
	d.output.Write([]byte("[::b]AVAILABLE COMMANDS[::b]\n\n"))

	// Group commands by category
	categories := map[string][]string{
		"Basic":     {"help", "status", "quit"},
		"Resources": {"gather", "resources"},
		"Villagers": {"recruit", "assign", "villagers"},
		"Buildings": {"build", "buildings"},
		"Research":  {"research", "technologies"},
		"Game":      {"save", "load"},
	}

	// Assign commands to proper categories
	commandCategories := make(map[string]map[string]string)
	for category, cmds := range categories {
		commandCategories[category] = make(map[string]string)
		for _, cmd := range cmds {
			if desc, ok := commands[cmd]; ok {
				commandCategories[category][cmd] = desc
			}
		}
	}

	// Add uncategorized commands
	commandCategories["Other"] = make(map[string]string)
	for cmd, desc := range commands {
		found := false
		for _, cmds := range categories {
			for _, c := range cmds {
				if c == cmd {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			commandCategories["Other"][cmd] = desc
		}
	}

	// Display commands by category
	for category, cmds := range commandCategories {
		if len(cmds) > 0 {
			d.output.Write([]byte(fmt.Sprintf("[#3498db]%s Commands:[#ffffff]\n", category)))
			for cmd, desc := range cmds {
				line := fmt.Sprintf("  [#f1c40f]%s[#ffffff]: %s\n", cmd, desc)
				d.output.Write([]byte(line))
			}
			d.output.Write([]byte("\n"))
		}
	}

	// Add keyboard shortcuts
	d.output.Write([]byte("[#3498db]Keyboard Shortcuts:[#ffffff]\n"))
	d.output.Write([]byte("  [#f1c40f]ESC[#ffffff]: Return focus to command input\n"))
	d.output.Write([]byte("  [#f1c40f]F1[#ffffff]: Show this help screen\n"))
}

// ShowMessage displays a message to the player
func (d *Display) ShowMessage(message string, style string) {
	var color string
	var prefix string

	switch style {
	case "info":
		color = "#3498db" // Blue
		prefix = "[i] "
	case "success":
		color = "#2ecc71" // Green
		prefix = "[+] "
	case "warning":
		color = "#f1c40f" // Yellow
		prefix = "[!] "
	case "error":
		color = "#e74c3c" // Red
		prefix = "[x] "
	case "highlight":
		color = "#9b59b6" // Purple
		prefix = "[*] "
	default:
		color = "#ffffff" // White
		prefix = ""
	}

	formattedMessage := fmt.Sprintf("[%s]%s%s[#ffffff]", color, prefix, message)
	d.output.Write([]byte(formattedMessage + "\n"))

	// Auto-scroll to the bottom
	d.output.ScrollToEnd()
}

// ShowAgeAdvancement displays age advancement notification
func (d *Display) ShowAgeAdvancement(newAge string) {
	// Create a themed modal for age advancement
	ageColors := map[string]string{
		"Stone Age":      "#7f8c8d", // Gray
		"Bronze Age":     "#d35400", // Bronze/Orange
		"Iron Age":       "#7f8c8d", // Gray-silver
		"Classical Age":  "#f1c40f", // Gold
		"Medieval Age":   "#8e44ad", // Royal purple
		"Renaissance":    "#3498db", // Blue
		"Industrial Age": "#e74c3c", // Red
		"Modern Age":     "#2ecc71", // Green
	}

	ageColor := "#2ecc71" // Default color
	if color, ok := ageColors[newAge]; ok {
		ageColor = color
	}

	messageText := fmt.Sprintf(`
[::b][#ffffff]CIVILIZATION ADVANCEMENT[::b]

Congratulations! Your civilization has advanced to the 
[::b][%s]%s[#ffffff][::b]

New buildings, technologies, and opportunities await you.
`, ageColor, newAge)

	modal := tview.NewModal().
		SetText(messageText).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Continue"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			d.pages.RemovePage("advancement")
			d.app.SetFocus(d.input)
		})

	// Add the modal as a page
	d.pages.AddPage("advancement", modal, true, true)
	d.app.SetFocus(modal)

	// Let the modal be visible for a moment
	time.Sleep(1 * time.Second)
}

// DisplayDashboard updates the dashboard with current game state
func (d *Display) DisplayDashboard(gameState game.GameState) {
	// Update status info
	d.status.Clear()
	d.status.Write([]byte(fmt.Sprintf("[#3498db]Age:[#ffffff] %s\n", gameState.Age)))
	d.status.Write([]byte(fmt.Sprintf("[#3498db]Tick:[#ffffff] %d\n", gameState.Tick)))

	// Calculate total villagers and display housing capacity
	totalVillagers := 0
	for _, v := range gameState.Villagers {
		totalVillagers += v.Count
	}
	capacityColor := "#2ecc71" // Green by default
	if totalVillagers >= gameState.VillagerCap {
		capacityColor = "#e74c3c" // Red if at capacity
	} else if float64(totalVillagers) >= float64(gameState.VillagerCap)*0.8 {
		capacityColor = "#f1c40f" // Yellow if close to capacity
	}

	d.status.Write([]byte(fmt.Sprintf("[#3498db]Housing:[#ffffff] %d/[%s]%d[#ffffff]\n",
		totalVillagers, capacityColor, gameState.VillagerCap)))

	// Show housing buildings available in the current age
	housingBuildings := []string{"hut"} // Only huts provide housing in this game
	for _, building := range housingBuildings {
		if count, exists := gameState.Buildings[building]; exists {
			d.status.Write([]byte(fmt.Sprintf("  [#3498db]%s:[#ffffff] %d\n",
				formatBuildingName(building), count)))
		}
	}

	// Update research display
	d.research.Clear()
	if gameState.Research.Current != "" {
		progressPercent := (gameState.Research.Progress / gameState.Research.Cost) * 100
		progressBar := createProgressBar(progressPercent, 20)

		d.research.Write([]byte(fmt.Sprintf("[#3498db]Researching:[#ffffff] %s\n", gameState.Research.Current)))
		d.research.Write([]byte(fmt.Sprintf("[#3498db]Progress:[#ffffff] %s %.1f%%\n",
			progressBar, progressPercent)))

		// Show completed research if any
		if len(gameState.Research.Researched) > 0 {
			d.research.Write([]byte(fmt.Sprintf("[#3498db]Completed:[#ffffff] %d technologies\n",
				len(gameState.Research.Researched))))
		}
	} else {
		d.research.Write([]byte("[#f1c40f]No active research[#ffffff]\n"))
		d.research.Write([]byte("Use 'research <technology>' command to start researching\n"))
	}

	// Update resources table with colors based on amounts
	d.resources.Clear()
	d.resources.SetCell(0, 0, &tview.TableCell{
		Text:       "Resource",
		Align:      tview.AlignLeft,
		Color:      tcell.ColorTeal,
		Attributes: tcell.AttrBold,
	})
	d.resources.SetCell(0, 1, &tview.TableCell{
		Text:       "Amount",
		Align:      tview.AlignRight,
		Color:      tcell.ColorGreen,
		Attributes: tcell.AttrBold,
	})

	// Add actual resources
	i := 1
	for resource, amount := range gameState.Resources {
		// Choose color based on amount
		resourceColor := tcell.ColorWhite
		if amount <= 5 {
			resourceColor = tcell.ColorRed
		} else if amount <= 20 {
			resourceColor = tcell.ColorYellow
		} else {
			resourceColor = tcell.ColorGreen
		}

		d.resources.SetCell(i, 0, &tview.TableCell{
			Text:  resource,
			Align: tview.AlignLeft,
		})
		d.resources.SetCell(i, 1, &tview.TableCell{
			Text:  strconv.FormatFloat(amount, 'f', 1, 64),
			Align: tview.AlignRight,
			Color: resourceColor,
		})
		i++
	}

	// Update villagers table with improved formatting
	d.villagers.Clear()
	d.villagers.SetCell(0, 0, &tview.TableCell{
		Text:       "Type",
		Align:      tview.AlignLeft,
		Color:      tcell.ColorTeal,
		Attributes: tcell.AttrBold,
	})
	d.villagers.SetCell(0, 1, &tview.TableCell{
		Text:       "Count",
		Align:      tview.AlignRight,
		Color:      tcell.ColorGreen,
		Attributes: tcell.AttrBold,
	})
	d.villagers.SetCell(0, 2, &tview.TableCell{
		Text:       "Assignment",
		Align:      tview.AlignLeft,
		Color:      tcell.ColorYellow,
		Attributes: tcell.AttrBold,
	})

	// Add actual villagers
	i = 1
	for vtype, info := range gameState.Villagers {
		if info.Count > 0 {
			d.villagers.SetCell(i, 0, &tview.TableCell{
				Text:  vtype,
				Align: tview.AlignLeft,
			})
			d.villagers.SetCell(i, 1, &tview.TableCell{
				Text:  strconv.Itoa(info.Count),
				Align: tview.AlignRight,
			})

			// Format assignments with improved readability
			var idleCount int
			assignmentTexts := make([]string, 0)

			for resource, count := range info.Assignment {
				if resource == "idle" {
					idleCount = count
					continue
				}

				if count > 0 {
					assignmentTexts = append(assignmentTexts,
						fmt.Sprintf("%s: %d", resource, count))
				}
			}
			assignments := ""
			if len(assignmentTexts) > 0 {
				assignments = joinWithCommas(assignmentTexts)
			}

			// Add idle count if there are idle villagers
			if idleCount > 0 {
				idleText := fmt.Sprintf("Idle: %d", idleCount)
				if assignments != "" {
					assignments += ", " + idleText
				} else {
					assignments = idleText
				}
			}

			// Add assignments to the table
			cellColor := tcell.ColorWhite
			if idleCount > 0 && idleCount == info.Count {
				cellColor = tcell.ColorYellow // All idle
			}

			d.villagers.SetCell(i, 2, &tview.TableCell{
				Text:  assignments,
				Align: tview.AlignLeft,
				Color: cellColor,
			})
			i++
		}
	}

	// Update buildings table with improved formatting
	d.buildings.Clear()
	d.buildings.SetCell(0, 0, &tview.TableCell{
		Text:       "Building",
		Align:      tview.AlignLeft,
		Color:      tcell.ColorTeal,
		Attributes: tcell.AttrBold,
	})
	d.buildings.SetCell(0, 1, &tview.TableCell{
		Text:       "Count",
		Align:      tview.AlignRight,
		Color:      tcell.ColorGreen,
		Attributes: tcell.AttrBold,
	})

	// Add actual buildings
	i = 1
	for building, count := range gameState.Buildings {
		if count > 0 {
			d.buildings.SetCell(i, 0, &tview.TableCell{
				Text:  building,
				Align: tview.AlignLeft,
			})
			d.buildings.SetCell(i, 1, &tview.TableCell{
				Text:  strconv.Itoa(count),
				Align: tview.AlignRight,
			})
			i++
		}
	}
}

// GetInput gets input from the user
func (d *Display) GetInput() (string, error) {
	// Focus on the input field
	d.app.SetFocus(d.input)

	// Wait for input from the channel
	return <-d.inputChan, nil
}

// Stop stops the tview application and resets the terminal
func (d *Display) Stop() {
	// Stop the application and clean up the terminal
	d.app.Stop()
}

// ShowLibraryContent displays library content in the output panel (CLI style)
func (d *Display) ShowLibraryContent(title, content string) {
	d.output.Clear()
	d.output.SetTextColor(tcell.ColorTeal)
	d.output.Write([]byte(fmt.Sprintf("[::b]%s[::b]\n\n", title)))
	d.output.SetTextColor(tcell.ColorWhite)
	d.output.Write([]byte(content + "\n"))
	d.output.Write([]byte("\n[#f1c40f]Type 'library' to see all topics, or 'library <topic>' to view another topic.[#ffffff]\n"))
	// Always restore input field focus
	d.app.SetFocus(d.input)
}

// ShowLibraryTopicsList displays the list of library topics in the output panel (CLI style)
func (d *Display) ShowLibraryTopicsList(topics map[string]string) {
	d.output.Clear()
	d.output.SetTextColor(tcell.ColorTeal)
	d.output.Write([]byte("[::b]LIBRARY TOPICS[::b]\n\n"))
	d.output.SetTextColor(tcell.ColorWhite)

	sortedIDs := make([]string, 0, len(topics))
	for id := range topics {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Strings(sortedIDs)

	for _, id := range sortedIDs {
		d.output.Write([]byte(fmt.Sprintf("  [#f1c40f]%s[#ffffff]: %s\n", id, topics[id])))
	}
	d.output.Write([]byte("\n[#f1c40f]Type 'library <topic>' to view details about a topic.[#ffffff]\n"))
	d.app.SetFocus(d.input)
}

// Helper function to create a text-based progress bar
func createProgressBar(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	completedWidth := int(float64(width) * percent / 100)

	bar := "["
	for i := 0; i < width; i++ {
		if i < completedWidth {
			bar += "#"
		} else {
			bar += "-"
		}
	}
	bar += "]"

	return bar
}

// Helper function to join strings with commas
func joinWithCommas(items []string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}

// Helper function to format building name for display
func formatBuildingName(name string) string {
	// Replace underscores with spaces
	formatted := strings.Replace(name, "_", " ", -1)
	// Capitalize first letter
	if len(formatted) > 0 {
		formatted = strings.ToUpper(formatted[:1]) + formatted[1:]
	}
	return formatted
}
