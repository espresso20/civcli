package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Settings provides the game settings interface
type Settings struct {
	ui         *UIManager
	view       *tview.Flex
	info       *tview.TextView
	menu       *tview.List
	returnPage string
}

// NewSettings creates a new settings screen
func NewSettings(ui *UIManager) *Settings {
	s := &Settings{
		ui:         ui,
		view:       tview.NewFlex(),
		info:       tview.NewTextView(),
		menu:       tview.NewList(),
		returnPage: "splash",
	}

	s.setupInfo()
	s.setupMenu()
	s.setupLayout()

	return s
}

// setupInfo creates the settings information display
func (s *Settings) setupInfo() {
	theme := s.ui.GetTheme()

	s.info.SetBorder(true).
		SetTitle(" ⚙️ Game Information ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	s.info.SetDynamicColors(true).
		SetWordWrap(true)

	// Set initial content without drawing
	content := s.generateInfoContent()
	s.info.SetText(content)
}

// setupMenu creates the settings menu
func (s *Settings) setupMenu() {
	theme := s.ui.GetTheme()

	s.menu.SetBorder(true).
		SetTitle(" 🔧 Settings Options ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)

	// Add menu options
	s.menu.AddItem("🔄 Refresh Info", "Update the displayed information", 'r', func() {
		s.updateInfo()
	})

	s.menu.AddItem("📁 Open Save Directory", "View save game location", 's', func() {
		s.ui.ShowMessage("Save directory: ./data/saves/", "info")
	})

	s.menu.AddItem("🔙 Back", "Return to previous screen", 'b', func() {
		s.ui.HideSettings()
	})

	// Set up input capture for navigation
	s.menu.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			s.ui.HideSettings()
			return nil
		}
		return event
	})
}

// setupLayout arranges the settings components
func (s *Settings) setupLayout() {
	s.view.SetDirection(tview.FlexColumn)

	s.view.
		AddItem(s.info, 0, 2, false). // Info panel (2/3 width)
		AddItem(s.menu, 25, 0, true)  // Menu panel (fixed width)
}

// updateInfo refreshes the information display
func (s *Settings) updateInfo() {
	content := s.generateInfoContent()
	s.info.SetText(content)
}

// generateInfoContent creates the settings information text
func (s *Settings) generateInfoContent() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "Unknown"
	}

	// Get save directory info
	saveDir := filepath.Join(cwd, "data", "saves")
	saveDirExists := "❌ Not found"
	if _, err := os.Stat(saveDir); err == nil {
		saveDirExists = "✅ Exists"
	}

	// Count save files
	saveCount := 0
	if files, err := os.ReadDir(saveDir); err == nil {
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
				saveCount++
			}
		}
	}

	// Get executable info
	execPath, err := os.Executable()
	if err != nil {
		execPath = "Unknown"
	}

	content := fmt.Sprintf(`[yellow::b]⚙️ CIV IDLE Settings & Information[white::-]

[cyan::b]📋 Game Version Information[white::-]
• Game Name: CIV IDLE
• Version: 1.0.0 (Development)
• Build Date: %s
• Platform: Terminal/Console

[cyan::b]📁 File System Information[white::-]
• Working Directory: %s
• Executable Path: %s
• Save Directory: %s
• Save Directory Status: %s
• Save Files Found: %d

[cyan::b]🎮 Game Data[white::-]
• Configuration: Default settings
• Theme: Modern Blue
• Input Method: Keyboard navigation
• Display Engine: tview terminal UI

[cyan::b]💾 Save Game Information[white::-]
• Save Format: JSON
• Auto-save: Enabled
• Save Location: ./data/saves/
• Backup System: Not implemented

[cyan::b]🔧 System Information[white::-]
• Terminal Support: ✅ Full color
• Unicode Support: ✅ Emojis enabled  
• Mouse Support: ✅ Available
• Resize Support: ✅ Dynamic layout

[green::b]💡 Tips & Notes[white::-]
• Save files are stored in JSON format for easy backup
• The game auto-saves progress at regular intervals
• Settings changes take effect immediately
• Use F1 for quick help from any screen
• Press Ctrl+Q to quit from anywhere

[yellow::b]🚀 Getting Started[white::-]
If this is your first time playing:
1. Return to the main menu
2. Select "Start New Civilization"
3. Press F1 anytime for help
4. Check the help system for detailed guides

Ready to build your civilization!`,
		time.Now().Format("2006-01-02"),
		cwd,
		execPath,
		saveDir,
		saveDirExists,
		saveCount)

	return content
}

// GetView returns the settings view
func (s *Settings) GetView() tview.Primitive {
	return s.view
}

// Focus sets focus to the settings screen
func (s *Settings) Focus() {
	s.ui.GetApp().SetFocus(s.menu)
}

// SetReturnPage sets which page to return to when settings is closed
func (s *Settings) SetReturnPage(page string) {
	s.returnPage = page
}

// GetReturnPage returns the page to return to when settings is closed
func (s *Settings) GetReturnPage() string {
	return s.returnPage
}
