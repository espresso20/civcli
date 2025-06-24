package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/user/civcli/game"
)

// DummyDisplay implements DisplayInterface for testing
type DummyDisplay struct{}

func (d *DummyDisplay) DisplayDashboard(state game.GameState) {}
func (d *DummyDisplay) ShowMessage(message, msgType string)   {}
func (d *DummyDisplay) ShowHelp(commands map[string]string)   {}
func (d *DummyDisplay) ShowAgeAdvancement(newAge string)      {}
func (d *DummyDisplay) GetInput() (string, error)             { return "", fmt.Errorf("no input") }
func (d *DummyDisplay) Stop()                                 {}

func main() {
	app := tview.NewApplication()

	// Create a game engine
	engine := game.NewGameEngine(&DummyDisplay{})

	// Scan for save files
	saveDir := "./data/saves"
	files, err := os.ReadDir(saveDir)
	if err != nil {
		panic(err)
	}

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Load Game Test - Press Enter to load, ESC to quit")

	var saveFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			saveFiles = append(saveFiles, file.Name())
			displayName := fmt.Sprintf("💾 %s", file.Name())
			info, _ := file.Info()
			secondaryText := fmt.Sprintf("Size: %d bytes", info.Size())
			list.AddItem(displayName, secondaryText, 0, nil)
		}
	}

	// Set up the selection handler - this is where the problem likely occurs
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		fmt.Printf("=== SELECTION HANDLER CALLED ===\n")
		fmt.Printf("Index: %d\n", index)
		fmt.Printf("MainText: %s\n", mainText)
		fmt.Printf("SaveFiles length: %d\n", len(saveFiles))

		if index < 0 || index >= len(saveFiles) {
			fmt.Printf("ERROR: Invalid index!\n")
			return
		}

		fileName := strings.TrimSuffix(saveFiles[index], ".json")
		fmt.Printf("About to load file: %s\n", fileName)

		// This is the critical part - loading the game
		fmt.Printf("Calling engine.LoadGame...\n")
		err := engine.LoadGame(fileName)
		fmt.Printf("LoadGame returned with err: %v\n", err)

		if err != nil {
			fmt.Printf("ERROR: Load failed: %v\n", err)
			// Show error modal
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Failed to load: %v", err)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(list, true)
				})
			app.SetRoot(modal, true)
		} else {
			fmt.Printf("SUCCESS: Load succeeded!\n")
			// Show success modal
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Successfully loaded: %s", fileName)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(list, true)
				})
			app.SetRoot(modal, true)
		}
		fmt.Printf("=== SELECTION HANDLER FINISHED ===\n")
	})

	// Add escape to quit
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			fmt.Printf("ESC pressed, stopping app\n")
			app.Stop()
			return nil
		}
		return event
	})

	fmt.Printf("Starting tview app...\n")

	// Auto-trigger a selection after 2 seconds to test
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Printf("Auto-triggering selection of index 0...\n")
		app.QueueUpdateDraw(func() {
			// Simulate pressing Enter on the first item
			if len(saveFiles) > 0 {
				// Manually call the selected function to test
				list.SetCurrentItem(0)
				// This should trigger the SetSelectedFunc
				app.SetFocus(list)
			}
		})
	}()

	// Run the app
	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}
