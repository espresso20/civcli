package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowHelp displays a formatted help screen with commands and keyboard shortcuts
func ShowHelp(pages *tview.Pages, commands map[string]string, returnPage string) {
	// Call the more advanced version with a default application
	app := tview.NewApplication()

	// Use the more feature-rich version
	createHelpScreen(app, pages, commands, returnPage, false)
}

// ShowHelpWithApp displays a formatted help screen with commands and keyboard shortcuts
// and properly handles focus and app-level input
func ShowHelpWithApp(app *tview.Application, pages *tview.Pages, commands map[string]string, returnPage string) {
	// Create the dynamic help screen
	createHelpScreen(app, pages, commands, returnPage, true)
}

// formatCategoryContent creates content for a specific help category
func formatCategoryContent(category string, commands map[string]string) string {
	var builder strings.Builder

	// Add title based on category
	builder.WriteString(fmt.Sprintf("[yellow::b]%s[-::-]\n", strings.ToUpper(category)))
	builder.WriteString("[yellow]" + strings.Repeat("─", 40) + "[-]\n\n")

	if len(commands) == 0 {
		builder.WriteString("[red]No commands available in this category.[-]\n")
		return builder.String()
	}

	// Sort the commands for consistent display
	commandNames := make([]string, 0, len(commands))
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	// Format each command with color and alignment
	for _, name := range commandNames {
		desc := commands[name]
		builder.WriteString(fmt.Sprintf("[green::b]%s[-::-]:\n  %s\n\n", name, desc))
	}

	return builder.String()
}

// formatKeyboardShortcuts returns formatted keyboard shortcuts
func formatKeyboardShortcuts() string {
	var builder strings.Builder

	builder.WriteString("[yellow::b]KEYBOARD SHORTCUTS[-::-]\n")
	builder.WriteString("[yellow]" + strings.Repeat("─", 40) + "[-]\n\n")
	builder.WriteString("[green::b]Ctrl+Q[-::-]: Quit the game\n")
	builder.WriteString("[green::b]Ctrl+C[-::-]: Quit the game\n")
	builder.WriteString("[green::b]Tab[-::-]: Toggle focus between components\n")
	builder.WriteString("[green::b]ESC[-::-]: Close popups and return to main screen\n")
	builder.WriteString("[green::b]↑/↓[-::-]: Navigate in menus or scroll content\n")
	builder.WriteString("[green::b]Enter[-::-]: Select an option\n")

	return builder.String()
}

// formatGameTips returns formatted game tips
func formatGameTips() string {
	var builder strings.Builder

	builder.WriteString("[yellow::b]GAME TIPS[-::-]\n")
	builder.WriteString("[yellow]" + strings.Repeat("─", 40) + "[-]\n\n")
	builder.WriteString("• Type [green::b]help <command>[-::-] for more information about a specific command\n")
	builder.WriteString("• Villagers need food to survive - assign them to gather resources\n")
	builder.WriteString("• Research technologies to unlock new buildings and resources\n")
	builder.WriteString("• Balance resource gathering with research and building construction\n")
	builder.WriteString("• Make sure to have enough villagers gathering food\n")
	builder.WriteString("• Use the [green::b]build[-::-] command to construct new buildings\n")
	builder.WriteString("• Remember to check your research progress regularly\n")

	return builder.String()
}

// createHelpScreen is an internal function that creates a dynamic help screen with multiple pages
func createHelpScreen(app *tview.Application, pages *tview.Pages, commands map[string]string, returnPage string, restoreCapture bool) {
	// Store the original input capture if needed
	var originalInputCapture func(event *tcell.EventKey) *tcell.EventKey
	if restoreCapture {
		originalInputCapture = app.GetInputCapture()
	}

	// Categorize commands
	gameCommands := make(map[string]string)
	resourceCommands := make(map[string]string)
	villagerCommands := make(map[string]string)
	buildingCommands := make(map[string]string)
	researchCommands := make(map[string]string)
	miscCommands := make(map[string]string)

	for name, desc := range commands {
		switch {
		case strings.HasPrefix(name, "help") || strings.HasPrefix(name, "quit") || strings.HasPrefix(name, "save") || strings.HasPrefix(name, "load"):
			gameCommands[name] = desc
		case strings.Contains(name, "resource") || strings.Contains(name, "food") || strings.Contains(name, "wood") || strings.Contains(name, "stone") || name == "gather" || name == "stats":
			resourceCommands[name] = desc
		case strings.Contains(name, "villager") || name == "assign" || name == "unassign" || name == "recruit":
			villagerCommands[name] = desc
		case strings.Contains(name, "build") || strings.Contains(name, "upgrade"):
			buildingCommands[name] = desc
		case strings.Contains(name, "research") || strings.Contains(name, "technology"):
			researchCommands[name] = desc
		default:
			miscCommands[name] = desc
		}
	}

	// Define menu options
	type menuItem struct {
		title   string
		content func() string
		pageID  string
	}

	menuItems := []menuItem{
		{
			title:   "Game Commands",
			content: func() string { return formatCategoryContent("Game Commands", gameCommands) },
			pageID:  "game-commands",
		},
		{
			title: "Resource Commands",
			content: func() string {
				content := formatCategoryContent("Resource Commands", resourceCommands)
				return content + "\n[yellow]Resource commands help you gather and manage your civilization's resources.[-]"
			},
			pageID: "resource-commands",
		},
		{
			title: "Villager Commands",
			content: func() string {
				content := formatCategoryContent("Villager Commands", villagerCommands)
				return content + "\n[yellow]Villager commands help you manage your population and their tasks.[-]"
			},
			pageID: "villager-commands",
		},
		{
			title:   "Building Commands",
			content: func() string { return formatCategoryContent("Building Commands", buildingCommands) },
			pageID:  "building-commands",
		},
		{
			title:   "Research Commands",
			content: func() string { return formatCategoryContent("Research Commands", researchCommands) },
			pageID:  "research-commands",
		},
		{
			title:   "Other Commands",
			content: func() string { return formatCategoryContent("Other Commands", miscCommands) },
			pageID:  "other-commands",
		},
		{
			title:   "Keyboard Shortcuts",
			content: formatKeyboardShortcuts,
			pageID:  "keyboard-shortcuts",
		},
		{
			title:   "Game Tips",
			content: formatGameTips,
			pageID:  "game-tips",
		},
	}

	// Base page names
	basePageName := "help-screen"
	mainPageName := basePageName + "-main"

	// Create all content views and pages
	contentViews := make(map[string]*tview.TextView)
	backButtons := make(map[string]*tview.Button)

	// Create the main menu list
	menuList := tview.NewList()
	menuList.SetBorder(true)
	menuList.SetTitle("Help Topics")
	menuList.SetTitleAlign(tview.AlignCenter)

	// Function to restore original state and return to game
	exitHelpSystem := func() {
		if restoreCapture {
			app.SetInputCapture(originalInputCapture)
		}
		pages.SwitchToPage(returnPage)
	}

	// Add each menu option
	for i, item := range menuItems { // Create a local copy of the item for the closure
		menuItem := item
		pageID := basePageName + "-" + menuItem.pageID

		// Add the menu item
		menuList.AddItem(menuItem.title, "", rune('1'+i), func() {
			// Switch to the content page
			pages.SwitchToPage(pageID)
			// Focus on the content view
			app.SetFocus(contentViews[pageID])
		})

		// Create the content view
		contentView := tview.NewTextView().
			SetDynamicColors(true).
			SetWordWrap(true).
			SetScrollable(true).
			SetChangedFunc(func() {
				app.Draw()
			}).
			SetText(menuItem.content())
		contentViews[pageID] = contentView

		// Create a back button for this page
		backButton := tview.NewButton("Back").
			SetSelectedFunc(func() {
				pages.SwitchToPage(mainPageName)
				app.SetFocus(menuList)
			})
		backButtons[pageID] = backButton

		// Create button layout
		buttonFlex := tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(nil, 0, 1, false).
			AddItem(backButton, 10, 0, true).
			AddItem(nil, 0, 1, false)

		// Create content layout
		contentFlex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(contentView, 0, 1, true).
			AddItem(buttonFlex, 3, 0, false)

		// Create frame
		contentFrame := tview.NewFrame(contentFlex).
			SetBorders(2, 2, 2, 2, 4, 4).
			AddText(menuItem.title, true, tview.AlignCenter, tcell.ColorYellow).
			AddText("Press ESC to return to menu", false, tview.AlignCenter, tcell.ColorWhite)

		// Add the page (but don't show it yet)
		pages.AddPage(pageID, contentFrame, true, false)
	}

	// Create close button for main menu
	closeButton := tview.NewButton("Close").
		SetSelectedFunc(exitHelpSystem)

	// Create button bar for main menu
	mainButtonFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(closeButton, 10, 0, true).
		AddItem(nil, 0, 1, false)

	// Create instruction text
	instructionText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Use TAB to focus on menu. Use arrow keys to navigate and Enter to select a topic. ESC to go back, or TAB to Close and press Enter.[-]").
		SetTextAlign(tview.AlignCenter)

	// Create main layout
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(instructionText, 1, 0, false).
		AddItem(menuList, 0, 1, true).
		AddItem(mainButtonFlex, 3, 0, false)

	// Create main frame
	mainFrame := tview.NewFrame(mainFlex).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("Civ-Idle Help", true, tview.AlignCenter, tcell.ColorYellow).
		AddText("Press ESC to return to game", false, tview.AlignCenter, tcell.ColorWhite)

	// Add main page
	pages.AddPage(mainPageName, mainFrame, true, true)

	// Create global input handler for all help screens
	if restoreCapture {
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			currentPage, _ := pages.GetFrontPage()

			// Only process our keys if we're on a help page
			if strings.HasPrefix(currentPage, basePageName) {
				if event.Key() == tcell.KeyEscape {
					if currentPage == mainPageName {
						// If on main menu, exit help system
						exitHelpSystem()
					} else {
						// If on a content page, return to main menu
						pages.SwitchToPage(mainPageName)
						app.SetFocus(menuList)
					}
					return nil
				} else if event.Key() == tcell.KeyTab {
					// Toggle focus based on current page
					if currentPage == mainPageName {
						// On main menu, toggle between list and close button
						if app.GetFocus() == menuList {
							app.SetFocus(closeButton)
						} else {
							app.SetFocus(menuList)
						}
					} else {
						// On content page, toggle between content and back button
						contentView := contentViews[currentPage]
						backButton := backButtons[currentPage]

						if app.GetFocus() == contentView {
							app.SetFocus(backButton)
						} else {
							app.SetFocus(contentView)
						}
					}
					return nil
				}
			}
			return event
		})
	}

	// Set initial focus to the menu list
	app.SetFocus(menuList)
}
