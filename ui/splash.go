package ui

import (
	"github.com/rivo/tview"
)

// SplashScreen provides a modern welcome screen
type SplashScreen struct {
	ui   *UIManager
	view *tview.Flex
	logo *tview.TextView
	menu *tview.List
}

// NewSplashScreen creates a new splash screen
func NewSplashScreen(ui *UIManager) *SplashScreen {
	s := &SplashScreen{
		ui:   ui,
		view: tview.NewFlex(),
		logo: tview.NewTextView(),
		menu: tview.NewList(),
	}

	s.setupLogo()
	s.setupMenu()
	s.setupLayout()

	return s
}

// setupLogo creates the game logo and description
func (s *SplashScreen) setupLogo() {
	logo := `
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║         ██████╗██╗██╗   ██╗    ██╗██████╗ ██╗     ███████╗                  ║
║        ██╔════╝██║██║   ██║    ██║██╔══██╗██║     ██╔════╝                  ║
║        ██║     ██║██║   ██║    ██║██║  ██║██║     █████╗                    ║
║        ██║     ██║╚██╗ ██╔╝    ██║██║  ██║██║     ██╔══╝                    ║
║        ╚██████╗██║ ╚████╔╝     ██║██████╔╝███████╗███████╗                  ║
║         ╚═════╝╚═╝  ╚═══╝      ╚═╝╚═════╝ ╚══════╝╚══════╝                  ║
║                                                                               ║
║                    🏛️  BUILD • GROW • CONQUER • ADVANCE  🏛️                    ║
║                                                                               ║
║           A command-line civilization builder that grows with you             ║
║                                                                               ║
║               📖 Manage resources, research technologies                       ║
║               🏘️  Build cities, expand your empire                            ║
║               ⚔️  Strategic decisions shape your destiny                       ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝`

	description := `
[::b][#4A90E2]Welcome to CIV IDLE![#ffffff][::b]

Experience the evolution of human civilization from the Stone Age to the Digital Era.
Every decision shapes your civilization's destiny.

[#E2E2E2]🎮 Easy to learn, challenging to master
🌍 Rich progression through historical ages
⚡ Real-time strategic decision making
📊 Detailed resource and population management

[#F39C12]Ready to begin your journey through history?[#ffffff]`

	s.logo.
		SetText(logo + description).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(false)
}

// setupMenu creates the main menu options
func (s *SplashScreen) setupMenu() {
	s.menu.
		SetBorder(true).
		SetTitle(" 🎯 Choose Your Path ").
		SetTitleAlign(tview.AlignCenter)

	// Add menu options
	s.menu.AddItem("🚀 Start New Civilization", "Begin your journey from the Stone Age", 'n', func() {
		s.ui.ShowDashboard()
	})

	s.menu.AddItem("📁 Load Saved Game", "Continue an existing civilization", 'l', func() {
		s.ui.ShowLoadGame()
	})

	s.menu.AddItem("❓ Help & Tutorial", "Learn how to play the game", 'h', func() {
		s.ui.ShowHelpSystem()
	})

	s.menu.AddItem("⚙️  Settings", "Configure game preferences", 's', func() {
		s.ui.ShowSettings()
	})

	s.menu.AddItem("🚪 Exit Game", "Return to terminal", 'q', func() {
		s.ui.Stop()
	})
}

// setupLayout arranges the splash screen components
func (s *SplashScreen) setupLayout() {
	// Create vertical layout
	s.view.SetDirection(tview.FlexRow)

	// Add components with spacing
	s.view.
		AddItem(nil, 2, 0, false).    // Top padding
		AddItem(s.logo, 0, 1, false). // Logo and description
		AddItem(nil, 1, 0, false).    // Spacing
		AddItem(s.menu, 12, 0, true). // Menu
		AddItem(nil, 2, 0, false)     // Bottom padding

	// Create horizontal centering
	centered := tview.NewFlex().SetDirection(tview.FlexColumn)

	centered.
		AddItem(nil, 0, 1, false).   // Left padding
		AddItem(s.view, 0, 3, true). // Main content
		AddItem(nil, 0, 1, false)    // Right padding

	s.view = centered
}

// GetView returns the splash screen view
func (s *SplashScreen) GetView() tview.Primitive {
	return s.view
}

// Focus sets focus to the splash screen
func (s *SplashScreen) Focus() {
	s.ui.app.SetFocus(s.menu)
}
