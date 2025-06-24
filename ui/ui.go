package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/civcli/game"
)

// Theme defines the color scheme for the UI
type Theme struct {
	Primary    tcell.Color
	Secondary  tcell.Color
	Accent     tcell.Color
	Success    tcell.Color
	Warning    tcell.Color
	Error      tcell.Color
	Background tcell.Color
	Foreground tcell.Color
	Border     tcell.Color
	Highlight  tcell.Color
}

// DefaultTheme provides a modern, clean color scheme
var DefaultTheme = Theme{
	Primary:    tcell.ColorBlue,
	Secondary:  tcell.ColorGray,
	Accent:     tcell.ColorYellow,
	Success:    tcell.ColorGreen,
	Warning:    tcell.ColorOrange,
	Error:      tcell.ColorRed,
	Background: tcell.ColorBlack,
	Foreground: tcell.ColorWhite,
	Border:     tcell.ColorBlue,
	Highlight:  tcell.ColorLightBlue,
}

// UIManager is the main UI controller
type UIManager struct {
	app   *tview.Application
	theme Theme
	pages *tview.Pages

	// Core components
	splash    *SplashScreen
	dashboard *Dashboard
	help      *HelpSystem
	settings  *Settings
	loadGame  *LoadGame

	// Game engine reference
	gameEngine *game.GameEngine

	// State
	currentPage string
	mu          sync.RWMutex
	inputChan   chan string
	running     bool
}

// NewUIManager creates a new UI manager with modern styling
func NewUIManager() *UIManager {
	ui := &UIManager{
		app:       tview.NewApplication(),
		theme:     DefaultTheme,
		pages:     tview.NewPages(),
		inputChan: make(chan string, 100),
		running:   false,
	}

	// Initialize components
	ui.splash = NewSplashScreen(ui)
	ui.dashboard = NewDashboard(ui)
	ui.help = NewHelpSystem(ui)
	ui.settings = NewSettings(ui)
	ui.loadGame = NewLoadGame(ui)

	// Set up the application
	ui.setupApplication()

	return ui
}

// setupApplication configures the main application
func (ui *UIManager) setupApplication() {
	// Set application theme
	ui.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.SetStyle(tcell.StyleDefault.
			Background(ui.theme.Background).
			Foreground(ui.theme.Foreground))
		return false
	})

	// Add pages
	ui.pages.AddPage("splash", ui.splash.GetView(), true, true)
	ui.pages.AddPage("dashboard", ui.dashboard.GetView(), true, false)
	ui.pages.AddPage("help", ui.help.GetView(), true, false)
	ui.pages.AddPage("settings", ui.settings.GetView(), true, false)
	ui.pages.AddPage("loadgame", ui.loadGame.GetView(), true, false)

	// Set root
	ui.app.SetRoot(ui.pages, true)

	// Global key bindings
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return ui.handleGlobalKeys(event)
	})

	ui.currentPage = "splash"
}

// handleGlobalKeys processes global keyboard shortcuts
func (ui *UIManager) handleGlobalKeys(event *tcell.EventKey) *tcell.EventKey {
	switch {
	case event.Key() == tcell.KeyCtrlQ:
		ui.Stop()
		return nil
	case event.Key() == tcell.KeyCtrlH && ui.currentPage != "help":
		ui.ShowHelpSystem()
		return nil
	case event.Key() == tcell.KeyEscape && ui.currentPage == "help":
		ui.HideHelp()
		return nil
	case event.Key() == tcell.KeyEscape && ui.currentPage == "settings":
		ui.HideSettings()
		return nil
	case event.Key() == tcell.KeyEscape && ui.currentPage == "loadgame":
		ui.HideLoadGame()
		return nil
	case event.Key() == tcell.KeyF1:
		ui.ShowHelpSystem()
		return nil
	}
	return event
}

// Start runs the UI application
func (ui *UIManager) Start() error {
	ui.mu.Lock()
	ui.running = true
	ui.mu.Unlock()

	return ui.app.Run()
}

// Stop stops the UI application
func (ui *UIManager) Stop() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if ui.running {
		ui.running = false
		ui.app.Stop()
	}
}

// ShowSplash displays the splash screen
func (ui *UIManager) ShowSplash() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.currentPage = "splash"
	ui.pages.SwitchToPage("splash")
}

// ShowDashboard displays the main game dashboard
func (ui *UIManager) ShowDashboard() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.currentPage = "dashboard"
	ui.pages.SwitchToPage("dashboard")
	ui.dashboard.Focus()
}

// SwitchToDashboardFromLoadGame safely switches to dashboard from load game without deadlock
func (ui *UIManager) SwitchToDashboardFromLoadGame() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.currentPage = "dashboard"
	ui.pages.SwitchToPage("dashboard")
	ui.dashboard.Focus()
}

// ShowHelpSystem displays the help system (public API for UI navigation)
func (ui *UIManager) ShowHelpSystem() {
	ui.showHelpSystem()
}

// HideHelp returns to the previous page from help
func (ui *UIManager) HideHelp() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	returnPage := ui.help.GetReturnPage()
	ui.currentPage = returnPage
	ui.pages.SwitchToPage(returnPage)

	if returnPage == "dashboard" {
		ui.dashboard.Focus()
	}
}

// ShowSettings displays the settings screen
func (ui *UIManager) ShowSettings() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.settings.SetReturnPage(ui.currentPage)
	ui.currentPage = "settings"
	ui.pages.SwitchToPage("settings")
	ui.settings.Focus()
}

// HideSettings returns to the previous page from settings
func (ui *UIManager) HideSettings() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	returnPage := ui.settings.GetReturnPage()
	ui.currentPage = returnPage
	ui.pages.SwitchToPage(returnPage)

	if returnPage == "dashboard" {
		ui.dashboard.Focus()
	} else if returnPage == "splash" {
		ui.splash.Focus()
	}
}

// ShowLoadGame displays the load game screen
func (ui *UIManager) ShowLoadGame() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.loadGame.SetReturnPage(ui.currentPage)
	ui.currentPage = "loadgame"
	ui.pages.SwitchToPage("loadgame")
	ui.loadGame.Focus()
}

// HideLoadGame returns to the previous page from load game
func (ui *UIManager) HideLoadGame() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	returnPage := ui.loadGame.GetReturnPage()
	ui.currentPage = returnPage
	ui.pages.SwitchToPage(returnPage)

	if returnPage == "dashboard" {
		ui.dashboard.Focus()
	} else if returnPage == "splash" {
		ui.splash.Focus()
	}
}

// GetInput returns user input from the dashboard
func (ui *UIManager) GetInput() (string, error) {
	select {
	case cmd := <-ui.inputChan:
		return cmd, nil
	case <-time.After(50 * time.Millisecond):
		return "", fmt.Errorf("no input ready")
	}
}

// SendInput sends input to the UI system
func (ui *UIManager) SendInput(input string) {
	select {
	case ui.inputChan <- input:
	default:
		// Channel full, ignore
	}
}

// UpdateGameState updates the dashboard with new game state
func (ui *UIManager) UpdateGameState(state game.GameState) {
	ui.dashboard.UpdateState(state)
}

// ShowMessage displays a message in the dashboard
func (ui *UIManager) ShowMessage(message, msgType string) {
	ui.dashboard.ShowMessage(message, msgType)
}

// GetTheme returns the current UI theme
func (ui *UIManager) GetTheme() Theme {
	return ui.theme
}

// GetApp returns the tview application instance
func (ui *UIManager) GetApp() *tview.Application {
	return ui.app
}

// GetPages returns the tview pages instance
func (ui *UIManager) GetPages() *tview.Pages {
	return ui.pages
}

// SetGameEngine sets the game engine reference for the UI
func (ui *UIManager) SetGameEngine(engine *game.GameEngine) {
	ui.gameEngine = engine
}

// GetGameEngine returns the game engine reference
func (ui *UIManager) GetGameEngine() *game.GameEngine {
	return ui.gameEngine
}

// DisplayInterface implementation for game engine compatibility

// ShowHelp displays help with the given commands (DisplayInterface method)
func (ui *UIManager) ShowHelp(commands map[string]string) {
	// Call the existing ShowHelp method - commands parameter is for legacy compatibility
	ui.showHelpSystem()
}

// showHelpSystem shows the help system (internal method)
func (ui *UIManager) showHelpSystem() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	ui.help.SetReturnPage(ui.currentPage)
	ui.currentPage = "help"
	ui.pages.SwitchToPage("help")
	ui.help.Focus()
}

// ShowAgeAdvancement displays age advancement notification
func (ui *UIManager) ShowAgeAdvancement(newAge string) {
	ui.ShowMessage(fmt.Sprintf("🎉 Congratulations! Your civilization has advanced to the %s!", newAge), "success")
}

// DisplayDashboard updates the dashboard with new game state
func (ui *UIManager) DisplayDashboard(state game.GameState) {
	ui.UpdateGameState(state)
}
