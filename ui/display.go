package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	mutex      sync.Mutex // Mutex for thread safety

	// Add CommandHandler reference to Display struct
	CommandHandler *game.CommandHandler

	// Track tab completion state
	tabCompletionState struct {
		shownOptions      bool
		lastCommand       string
		selectedOptionIdx int
		currentOptions    []string
	}
}

// NewDisplay creates a new display
// Update NewDisplay to accept CommandHandler as a parameter
func NewDisplay(commandHandler *game.CommandHandler) *Display {
	d := &Display{
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		inputChan:      make(chan string),
		CommandHandler: commandHandler, // Assign CommandHandler
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
		SetFieldTextColor(tcell.ColorWhite).
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

	// Correct access to GetCommandList and implement manual autocomplete
	d.input.SetChangedFunc(func(text string) {
		if text == "" {
			return
		}

		// Filter commands based on current text
		var suggestions []string
		for _, cmd := range d.CommandHandler.GetCommandList() {
			if strings.HasPrefix(cmd, text) {
				suggestions = append(suggestions, cmd)
			}
		}

		// Display suggestions
		// TODO: (this can be enhanced with a dropdown or similar UI)
		if len(suggestions) > 0 {
			d.output.Clear()
			for _, suggestion := range suggestions {
				d.output.Write([]byte(suggestion + "\n"))
			}
		}
	})

	// Enhance input field to support Minecraft-style tab completion
	d.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			currentText := d.input.GetText()

			// If text is empty, suggest the first command
			if currentText == "" {
				commands := d.CommandHandler.GetCommandList()
				if len(commands) > 0 {
					d.input.SetText(commands[0])
				}
				return nil
			}

			// Split the text into parts
			parts := strings.Fields(currentText)

			// Determine what we're completing (a command or a subcommand)
			if len(parts) == 0 {
				// This shouldn't happen if currentText isn't empty, but just in case
				return event
			}

			// Get the part we're currently typing (last word)
			lastWord := parts[len(parts)-1]

			// Check if we just finished typing a complete word (space at the end)
			isCompletedWord := strings.HasSuffix(currentText, " ")

			// Determine suggestions based on context
			var suggestions []string

			if len(parts) == 1 && !isCompletedWord {
				// We're typing the first word (a command) - suggest from command list
				commands := d.CommandHandler.GetCommandList()
				for _, cmd := range commands {
					if strings.HasPrefix(cmd, lastWord) {
						suggestions = append(suggestions, cmd)
					}
				}
			} else {
				// We're typing a subcommand or have completed the command
				command := parts[0]

				// Get suggestions based on the command
				switch command {
				case "gather":
					suggestions = []string{"food", "wood", "stone"}
				case "build":
					suggestions = []string{"hut", "farm", "mine"}
				case "research":
					suggestions = []string{"agriculture", "mining", "writing"}
				case "assign":
					suggestions = []string{"villager", "task", "building"}
				case "library":
					suggestions = []string{"villagers", "resources", "buildings", "ages", "commands", "tips"}
				default:
					suggestions = []string{"help", "status", "quit"}
				}

				// If we're not at a space, filter based on what we've typed so far
				if !isCompletedWord && len(parts) > 1 {
					filteredSuggestions := []string{}
					for _, suggestion := range suggestions {
						if strings.HasPrefix(suggestion, lastWord) {
							filteredSuggestions = append(filteredSuggestions, suggestion)
						}
					}
					suggestions = filteredSuggestions
				}
			}

			// If we have suggestions, use them
			if len(suggestions) > 0 {
				// Check if we're showing options for a command
				if isCompletedWord && len(parts) == 1 {
					// Static variable to track if we've shown the options
					// It will persist between function calls
					static := struct {
						optionsShown bool
						currentIdx   int
						lastCommand  string
					}{}

					command := parts[0]

					// Reset state if the command changed
					if static.lastCommand != command {
						static.optionsShown = false
						static.currentIdx = 0
						static.lastCommand = command
					}

					// If options haven't been shown yet, show them first
					if !static.optionsShown {
						d.output.Clear()
						d.output.SetTextColor(tcell.ColorTeal)
						d.output.Write([]byte(fmt.Sprintf("[::b]Available options for '%s':[::b]\n\n", command)))
						d.output.SetTextColor(tcell.ColorYellow)

						for _, suggestion := range suggestions {
							d.output.Write([]byte("  " + suggestion + "\n"))
						}

						d.output.Write([]byte("\n[#f1c40f]Press Tab again to cycle through options[#ffffff]\n"))

						// Mark that we've shown options
						static.optionsShown = true
						return nil
					} else {
						// Options have been shown, now we cycle through them
						d.input.SetText(currentText + suggestions[static.currentIdx])
						static.currentIdx = (static.currentIdx + 1) % len(suggestions)
						return nil
					}
				} else if isCompletedWord && len(parts) > 1 {
					// Cycling through subcommand options after a command and its first subcommand
					// Get the current suggestion to append
					currentIndex := 0
					if len(parts) > 1 && parts[len(parts)-1] != "" {
						// We already have a subcommand, find its index
						for i, suggestion := range suggestions {
							if parts[len(parts)-1] == suggestion {
								currentIndex = (i + 1) % len(suggestions)
								break
							}
						}
					}

					// Replace the last part with the next suggestion
					if len(parts) > 1 && parts[len(parts)-1] != "" {
						parts[len(parts)-1] = suggestions[currentIndex]
						d.input.SetText(strings.Join(parts, " ") + " ")
					} else {
						// Append the first suggestion
						d.input.SetText(currentText + suggestions[currentIndex])
					}
					return nil
				}

				// We're completing a partial word
				// Try to find if we're cycling through suggestions
				currentSuggestionIndex := -1
				for i, suggestion := range suggestions {
					if lastWord == suggestion {
						currentSuggestionIndex = i
						break
					}
				}

				// Get the next suggestion (or the first one if not found)
				nextIndex := (currentSuggestionIndex + 1) % len(suggestions)
				chosenSuggestion := suggestions[nextIndex]

				// Replace the last word with the suggestion
				if len(parts) == 1 {
					// For the first word, just set it directly with a space
					d.input.SetText(chosenSuggestion + " ")
				} else {
					// For subcommands, preserve the command and replace the last word
					newParts := append(parts[:len(parts)-1], chosenSuggestion)
					d.input.SetText(strings.Join(newParts, " ") + " ")
				}
			}

			return nil
		}
		return event
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

	// topDashboard := tview.NewFlex().
	// 	SetDirection(tview.FlexColumn).
	// 	AddItem(leftPanel, 0, 1, false).
	// 	AddItem(d.rightPanel, 0, 1, false)

	// Add banner panel to the top of the display
	banner := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	banner.SetBorder(true).SetTitle("")

	// Banner content: Game name centered, version lower left, creator lower right
	bannerText := `[::b][green]CivIdleCli - A Command Line Civilization Builder[::b]
` +
		`[white][v1.0.0][::d]` + strings.Repeat(" ", 60) + `[white]Creator: Espresso[::d]`
	banner.SetText(bannerText)

	// Add the banner to the top of the dashboard
	// Adjust dashboard layout to include the banner at the top
	d.dashboard = tview.NewGrid().
		SetRows(5, 20, 0).
		SetColumns(0).
		AddItem(banner, 0, 0, 1, 2, 0, 0, false).
		AddItem(leftPanel, 1, 0, 1, 1, 0, 0, false).
		AddItem(d.rightPanel, 1, 1, 1, 1, 0, 0, false).
		AddItem(d.output, 2, 0, 1, 2, 0, 0, false)

	// Create main layout
	d.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(d.dashboard, 0, 1, false).
		AddItem(d.input, 1, 1, true)

	// Create intro page with improved styling
	introText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetScrollable(false) // Disable scrolling to show all content at once
	introText.SetText(`
[::b][green]CivIdleCli - A Command Line Civilization Builder Game[::d][::b]

[blue]Welcome to CivIdleCli![::d] You are the leader of a small tribe in the Stone Age.

[yellow]Your Mission:[::d]
- Gather resources to survive and grow
- Build structures to advance your civilization
- Research technologies to unlock new possibilities
- Assign villagers efficiently to different tasks
- Guide your people through the ages of history

[red]Type 'help' at any time to see available commands.[::d]

[red]You can read the wiki library for further knowledge on how to play.[::d]

[::b][green]Press Enter to begin your journey...[::d][::b]
	`)

	introPage := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(introText, 20, 1, false).
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
			d.output.Write([]byte(fmt.Sprintf("[blue]%s Commands:[::d]\n", category)))
			for cmd, desc := range cmds {
				line := fmt.Sprintf("  [yellow]%s[::d]: %s\n", cmd, desc)
				d.output.Write([]byte(line))
			}
			d.output.Write([]byte("\n"))
		}
	}

	// Add keyboard shortcuts
	d.output.Write([]byte("[blue]Keyboard Shortcuts:[::d]\n"))
	d.output.Write([]byte("  [yellow]ESC[::d]: Return focus to command input\n"))
	d.output.Write([]byte("  [yellow]F1[::d]: Show this help screen\n"))
}

// ShowMessage displays a message to the player
func (d *Display) ShowMessage(message string, style string) {
	var color string
	var prefix string

	switch style {
	case "info":
		color = "blue" // Blue
		prefix = "[i] "
	case "success":
		color = "green" // Green
		prefix = "[+] "
	case "warning":
		color = "yellow" // Yellow
		prefix = "[!] "
	case "error":
		color = "red" // Red
		prefix = "[x] "
	case "highlight":
		color = "purple" // Purple
		prefix = "[*] "
	default:
		color = "::d" // Default color
		prefix = ""
	}

	formattedMessage := fmt.Sprintf("[%s]%s%s[::d]", color, prefix, message)
	d.output.Write([]byte(formattedMessage + "\n"))

	// Auto-scroll to the bottom
	d.output.ScrollToEnd()
}

// ShowAgeAdvancement displays age advancement notification
func (d *Display) ShowAgeAdvancement(newAge string) {
	// Create a themed modal for age advancement
	ageColors := map[string]string{
		"Stone Age":      "gray",    // Gray
		"Bronze Age":     "orange",  // Bronze/Orange
		"Iron Age":       "gray",    // Gray-silver
		"Classical Age":  "yellow",  // Gold
		"Medieval Age":   "purple",  // Royal purple
		"Renaissance":    "blue",    // Blue
		"Industrial Age": "red",     // Red
		"Modern Age":     "green",   // Green
	}

	ageColor := "green" // Default color
	if color, ok := ageColors[newAge]; ok {
		ageColor = color
	}

	messageText := fmt.Sprintf(`
[::b][::d]CIVILIZATION ADVANCEMENT[::b]

Congratulations! Your civilization has advanced to the 
[::b][%s]%s[::d][::b]

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
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Update status info
	d.status.Clear()
	d.status.Write([]byte(fmt.Sprintf("[blue]Age:[::d] %s\n", gameState.Age)))
	d.status.Write([]byte(fmt.Sprintf("[blue]Tick:[::d] %d\n", gameState.Tick)))
	d.status.Write([]byte(fmt.Sprintf("[blue]Tick Rate:[::d] %.1fs per tick\n", gameState.TickDurationSeconds)))

	// Force draw to ensure status gets updated
	d.app.Draw()

	// Calculate total villagers and display housing capacity
	totalVillagers := 0
	for _, v := range gameState.Villagers {
		totalVillagers += v.Count
	}
	capacityColor := "green" // Green by default
	if totalVillagers >= gameState.VillagerCap {
		capacityColor = "red" // Red if at capacity
	} else if float64(totalVillagers) >= float64(gameState.VillagerCap)*0.8 {
		capacityColor = "yellow" // Yellow if close to capacity
	}

	d.status.Write([]byte(fmt.Sprintf("[blue]Housing:[::d] %d/[%s]%d[::d]\n",
		totalVillagers, capacityColor, gameState.VillagerCap)))

	// Show housing buildings available in the current age
	housingBuildings := []string{"hut"} // Only huts provide housing in this game
	for _, building := range housingBuildings {
		if count, exists := gameState.Buildings[building]; exists {
			d.status.Write([]byte(fmt.Sprintf("  [blue]%s:[::d] %d\n",
				formatBuildingName(building), count)))
		}
	}

	// Update research display
	d.research.Clear()
	if gameState.Research.Current != "" {
		progressPercent := (gameState.Research.Progress / gameState.Research.Cost) * 100
		progressBar := createProgressBar(progressPercent, 20)

		d.research.Write([]byte(fmt.Sprintf("[blue]Researching:[::d] %s\n", gameState.Research.Current)))
		d.research.Write([]byte(fmt.Sprintf("[blue]Progress:[::d] %s %.1f%%\n",
			progressBar, progressPercent)))

		// Show completed research if any
		if len(gameState.Research.Researched) > 0 {
			d.research.Write([]byte(fmt.Sprintf("[blue]Completed:[::d] %d technologies\n",
				len(gameState.Research.Researched))))
		}
	} else {
		d.research.Write([]byte("[yellow]No active research[::d]\n"))
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

	// Define a fixed order for resources
	resourceOrder := []string{
		"food",
		"wood",
		"stone",
		"gold",
		"knowledge",
		// Add any other resources that might be added later
	}

	// Add resources in the fixed order
	i := 1
	for _, resource := range resourceOrder {
		amount, exists := gameState.Resources[resource]
		if !exists {
			continue // Skip if the resource doesn't exist yet
		}

		// Choose color based on amount
		resourceColor := tcell.ColorDefault
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

	// Add any resources that weren't in our predefined list
	// This ensures we don't miss any new resources added in the future
	for resource, amount := range gameState.Resources {
		// Skip resources we've already displayed
		isKnownResource := false
		for _, knownResource := range resourceOrder {
			if resource == knownResource {
				isKnownResource = true
				break
			}
		}
		if isKnownResource {
			continue
		}

		// Choose color based on amount
		resourceColor := tcell.ColorDefault
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
			cellColor := tcell.ColorDefault
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
	d.output.Write([]byte("\n[yellow]Type 'library' to see all topics, or 'library <topic>' to view another topic.[::d]\n"))
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
		d.output.Write([]byte(fmt.Sprintf("  [yellow]%s[::d]: %s\n", id, topics[id])))
	}
	d.output.Write([]byte("\n[yellow]Type 'library <topic>' to view details about a topic.[::d]\n"))
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
