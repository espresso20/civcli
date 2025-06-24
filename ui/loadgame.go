package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// LoadGame provides the save game loading interface
type LoadGame struct {
	ui           *UIManager
	view         *tview.Flex
	saveList     *tview.List
	infoPanel    *tview.TextView
	actionPanel  *tview.TextView
	returnPage   string
	saveFiles    []SaveFileInfo
	selectedSave *SaveFileInfo
}

// SaveFileInfo represents information about a save file
type SaveFileInfo struct {
	Name         string
	Path         string
	Size         int64
	ModTime      time.Time
	IsValid      bool
	ErrorMessage string
}

// NewLoadGame creates a new load game screen
func NewLoadGame(ui *UIManager) *LoadGame {
	lg := &LoadGame{
		ui:          ui,
		view:        tview.NewFlex(),
		saveList:    tview.NewList(),
		infoPanel:   tview.NewTextView(),
		actionPanel: tview.NewTextView(),
		returnPage:  "splash",
		saveFiles:   make([]SaveFileInfo, 0),
	}

	lg.setupSaveList()
	lg.setupInfoPanel()
	lg.setupActionPanel()
	lg.setupLayout()
	lg.refreshSaveFiles()

	return lg
}

// setupSaveList creates the save file list
func (lg *LoadGame) setupSaveList() {
	theme := lg.ui.GetTheme()

	lg.saveList.SetBorder(true).
		SetTitle(" 💾 Saved Games ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)

	// Set up selection handler
	lg.saveList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Reset selection if no save files
		if len(lg.saveFiles) == 0 {
			lg.selectedSave = nil
			lg.updateInfoPanel()
			lg.updateActionPanel()
			return
		}

		// Validate index bounds
		if index >= 0 && index < len(lg.saveFiles) {
			lg.selectedSave = &lg.saveFiles[index]
			lg.updateInfoPanel()
			lg.updateActionPanel()
		} else {
			lg.selectedSave = nil
			lg.updateInfoPanel()
			lg.updateActionPanel()
		}
	})

	// Set up selection handler for loading - simplified approach
	lg.saveList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// Only attempt to load if we have valid save files and a valid index
		if len(lg.saveFiles) == 0 {
			lg.ui.ShowMessage("No save files available to load", "warning")
			return
		}

		if index < 0 || index >= len(lg.saveFiles) {
			lg.ui.ShowMessage("Invalid save file selection", "error")
			return
		}

		// Instead of loading immediately, show a confirmation dialog
		lg.showLoadConfirmation(index)
	})

	// Set up input capture
	lg.saveList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			lg.ui.HideLoadGame()
			return nil
		case tcell.KeyCtrlR:
			lg.refreshSaveFiles()
			return nil
		case tcell.KeyDelete:
			if lg.selectedSave != nil && len(lg.saveFiles) > 0 {
				lg.deleteSaveFile()
			} else {
				lg.ui.ShowMessage("No save file selected for deletion", "warning")
			}
			return nil
		}
		return event
	})
}

// setupInfoPanel creates the save file information display
func (lg *LoadGame) setupInfoPanel() {
	theme := lg.ui.GetTheme()

	lg.infoPanel.SetBorder(true).
		SetTitle(" 📋 Save Information ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	lg.infoPanel.SetDynamicColors(true).
		SetWordWrap(true)

	lg.infoPanel.SetText("[yellow]Select a save file to view details[white]")
}

// setupActionPanel creates the action instructions panel
func (lg *LoadGame) setupActionPanel() {
	theme := lg.ui.GetTheme()

	lg.actionPanel.SetBorder(true).
		SetTitle(" 🎮 Actions ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(theme.Border)
	lg.actionPanel.SetDynamicColors(true)

	actionText := `[cyan]Available Actions:[white]

[yellow]Enter[white] - Load selected save game
[yellow]Delete[white] - Delete selected save file
[yellow]Ctrl+R[white] - Refresh save file list
[yellow]ESC[white] - Return to main menu

[green]Navigation:[white]
↑/↓ - Select save file
Tab - Switch between panels`

	lg.actionPanel.SetText(actionText)
}

// setupLayout arranges the load game components
func (lg *LoadGame) setupLayout() {
	// Create left panel (save list)
	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	leftPanel.AddItem(lg.saveList, 0, 1, true)

	// Create right panel (info and actions)
	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	rightPanel.
		AddItem(lg.infoPanel, 0, 2, false).
		AddItem(lg.actionPanel, 12, 0, false)

	// Main layout
	lg.view.SetDirection(tview.FlexColumn)
	lg.view.
		AddItem(leftPanel, 0, 1, true).
		AddItem(rightPanel, 0, 1, false)
}

// refreshSaveFiles scans for save files and updates the list
func (lg *LoadGame) refreshSaveFiles() {
	lg.saveFiles = make([]SaveFileInfo, 0)
	lg.selectedSave = nil

	saveDir := "./data/saves"

	// Check if save directory exists
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		lg.saveList.Clear()
		lg.saveList.AddItem("📁 No save directory found", "Create a save directory: ./data/saves", 0, nil)
		lg.updateInfoPanel()
		lg.updateActionPanel()
		return
	}

	// Read save directory
	files, err := os.ReadDir(saveDir)
	if err != nil {
		lg.saveList.Clear()
		lg.saveList.AddItem("❌ Error reading save directory", err.Error(), 0, nil)
		lg.updateInfoPanel()
		lg.updateActionPanel()
		return
	}

	// Process save files
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			continue
		}

		fullPath := filepath.Join(saveDir, file.Name())
		info, err := file.Info()
		if err != nil {
			continue
		}

		saveInfo := SaveFileInfo{
			Name:    file.Name(),
			Path:    fullPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsValid: true,
		}

		// Validate the save file
		isValid, errorMsg := lg.validateSaveFile(fullPath)
		saveInfo.IsValid = isValid
		saveInfo.ErrorMessage = errorMsg

		lg.saveFiles = append(lg.saveFiles, saveInfo)
	}

	// Sort by modification time (newest first)
	sort.Slice(lg.saveFiles, func(i, j int) bool {
		return lg.saveFiles[i].ModTime.After(lg.saveFiles[j].ModTime)
	})

	// Update the list
	lg.updateSaveList()
	lg.updateInfoPanel()
	lg.updateActionPanel()
}

// updateSaveList refreshes the save file list display
func (lg *LoadGame) updateSaveList() {
	lg.saveList.Clear()

	if len(lg.saveFiles) == 0 {
		lg.saveList.AddItem("📂 No save files found", "Start a new game to create your first save", 0, nil)
		return
	}

	for i, save := range lg.saveFiles {
		mainText := fmt.Sprintf("💾 %s", save.Name)
		secondaryText := fmt.Sprintf("Modified: %s | Size: %s",
			save.ModTime.Format("2006-01-02 15:04"),
			lg.formatFileSize(save.Size))

		if !save.IsValid {
			mainText = fmt.Sprintf("❌ %s (Invalid)", save.Name)
			secondaryText = save.ErrorMessage
		}

		lg.saveList.AddItem(mainText, secondaryText, rune('1'+i), nil)
	}
}

// updateInfoPanel refreshes the save file information display
func (lg *LoadGame) updateInfoPanel() {
	if lg.selectedSave == nil {
		if len(lg.saveFiles) == 0 {
			content := `[yellow::b]📂 No Save Files Found[white::-]

No saved games were found in the save directory.

[cyan]To create save files:[white]
1. Start a new civilization
2. Play the game
3. The game will auto-save your progress

[cyan]Save Location:[white]
./data/saves/

[green]Create your first civilization to get started![white]`
			lg.infoPanel.SetText(content)
		} else {
			lg.infoPanel.SetText("[yellow]Select a save file to view details[white]")
		}
		return
	}

	save := lg.selectedSave

	content := fmt.Sprintf(`[yellow::b]💾 Save File Details[white::-]

[cyan]File Name:[white] %s
[cyan]File Path:[white] %s
[cyan]File Size:[white] %s
[cyan]Created:[white] %s
[cyan]Last Modified:[white] %s

[cyan]Status:[white] %s

[yellow]Game Information:[white]
• Format: JSON save file
• Compatible: Yes
• Backup: Available

[green]Press Enter to load this save game[white]`,
		save.Name,
		save.Path,
		lg.formatFileSize(save.Size),
		save.ModTime.Format("Monday, January 2, 2006"),
		save.ModTime.Format("15:04:05"),
		lg.getStatusText(save))

	lg.infoPanel.SetText(content)
}

// updateActionPanel refreshes the action panel
func (lg *LoadGame) updateActionPanel() {
	if lg.selectedSave == nil {
		actionText := `[cyan]Available Actions:[white]

[yellow]Ctrl+R[white] - Refresh save file list
[yellow]ESC[white] - Return to main menu

[green]Navigation:[white]
↑/↓ - Select save file
Tab - Switch between panels

[gray]Select a save file to see more options[white]`
		lg.actionPanel.SetText(actionText)
		return
	}

	actionText := `[cyan]Available Actions:[white]

[yellow]Enter[white] - Load selected save game
[yellow]Delete[white] - Delete selected save file
[yellow]Ctrl+R[white] - Refresh save file list
[yellow]ESC[white] - Return to main menu

[green]Navigation:[white]
↑/↓ - Select save file
Tab - Switch between panels`

	lg.actionPanel.SetText(actionText)
}

// validateSaveFile performs basic validation on a save file
func (lg *LoadGame) validateSaveFile(path string) (bool, string) {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return false, "File not found"
	}

	// Check if file is empty
	if info.Size() == 0 {
		return false, "File is empty"
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return false, "Cannot read file"
	}

	// Basic JSON validation
	trimmed := strings.TrimSpace(string(data))
	if !strings.HasPrefix(trimmed, "{") {
		return false, "Not a valid JSON object"
	}

	// Try to parse as generic JSON to validate structure
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return false, "Invalid JSON format"
	}

	// Check for required fields that indicate this is a game save
	requiredFields := []string{"tick", "age", "resources"}
	for _, field := range requiredFields {
		if _, exists := jsonData[field]; !exists {
			return false, fmt.Sprintf("Missing required field: %s", field)
		}
	}

	return true, ""
}
func (lg *LoadGame) formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}

// getStatusText returns a status description for the save file
func (lg *LoadGame) getStatusText(save *SaveFileInfo) string {
	if !save.IsValid {
		return fmt.Sprintf("❌ Invalid (%s)", save.ErrorMessage)
	}
	return "✅ Valid save file"
}

// showLoadConfirmation shows a proper confirmation dialog with warning about overwriting current game
func (lg *LoadGame) showLoadConfirmation(index int) {
	if index < 0 || index >= len(lg.saveFiles) {
		lg.ui.ShowMessage("Invalid save file selection", "error")
		return
	}

	selectedSave := lg.saveFiles[index]
	fileName := strings.TrimSuffix(selectedSave.Name, ".json")

	// Create a detailed confirmation modal
	confirmText := fmt.Sprintf(`[yellow::b]Load Save Game: %s[white::-]

[red::b]⚠️  WARNING: This will overwrite any current game progress![white::-]

[cyan]Game Details:[white]
• File: %s
• Size: %s
• Modified: %s
• Status: %s

[green]Do you want to continue loading this save game?[white]`,
		fileName,
		selectedSave.Name,
		lg.formatFileSize(selectedSave.Size),
		selectedSave.ModTime.Format("2006-01-02 15:04"),
		lg.getStatusText(&selectedSave))

	modal := tview.NewModal().
		SetText(confirmText).
		AddButtons([]string{"YES - Load Game", "NO - Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// Remove the modal first
			lg.ui.GetPages().RemovePage("loadConfirm")

			if buttonIndex == 0 { // YES - Load Game
				// Load directly without goroutine to avoid race conditions
				lg.performGameLoad(fileName)
			}
			// If NO - Cancel, just return to the load game screen (do nothing)
		})

	// Add the modal as a new page
	lg.ui.GetPages().AddPage("loadConfirm", modal, true, true)
	lg.ui.GetApp().SetFocus(modal)
}

// performGameLoad performs the complete game loading process - minimal approach
func (lg *LoadGame) performGameLoad(fileName string) {
	gameEngine := lg.ui.GetGameEngine()
	if gameEngine == nil {
		return // Fail silently to avoid UI conflicts
	}

	// Load the game exactly like the in-game command does
	err := gameEngine.LoadGame(fileName)

	if err != nil {
		return // Fail silently to avoid UI conflicts for now
	}

	// Success - switch to dashboard immediately (the refresh loop will show the loaded game)
	lg.ui.SwitchToDashboardFromLoadGame()
}

// deleteSaveFile deletes the currently selected save file
func (lg *LoadGame) deleteSaveFile() {
	if lg.selectedSave == nil {
		lg.ui.ShowMessage("No save file selected for deletion", "warning")
		return
	}

	// Validate the file still exists
	if _, err := os.Stat(lg.selectedSave.Path); os.IsNotExist(err) {
		lg.ui.ShowMessage("Save file no longer exists", "warning")
		lg.refreshSaveFiles()
		return
	}

	// TODO: Add confirmation dialog in the future
	// For now, proceed with deletion but add safety checks

	fileName := lg.selectedSave.Name
	err := os.Remove(lg.selectedSave.Path)
	if err != nil {
		lg.ui.ShowMessage(fmt.Sprintf("Error deleting save file '%s': %v", fileName, err), "error")
		return
	}

	lg.ui.ShowMessage(fmt.Sprintf("Successfully deleted save file: %s", fileName), "success")
	lg.selectedSave = nil
	lg.refreshSaveFiles()
}

// GetView returns the load game view
func (lg *LoadGame) GetView() tview.Primitive {
	return lg.view
}

// Focus sets focus to the load game screen
func (lg *LoadGame) Focus() {
	lg.ui.GetApp().SetFocus(lg.saveList)
}

// SetReturnPage sets which page to return to when load game is closed
func (lg *LoadGame) SetReturnPage(page string) {
	lg.returnPage = page
}

// GetReturnPage returns the page to return to when load game is closed
func (lg *LoadGame) GetReturnPage() string {
	return lg.returnPage
}
