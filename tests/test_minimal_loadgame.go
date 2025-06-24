package main

import (
	"fmt"
	"os"
	"strings"

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
	list.SetBorder(true).SetTitle("Load Game Test")

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
		fmt.Printf("Selection handler called: index=%d, mainText=%s\n", index, mainText)

		if index < 0 || index >= len(saveFiles) {
			fmt.Printf("Invalid index!\n")
			return
		}

		fileName := strings.TrimSuffix(saveFiles[index], ".json")
		fmt.Printf("About to load: %s\n", fileName)

		// This is the critical part - loading the game
		err := engine.LoadGame(fileName)
		if err != nil {
			fmt.Printf("Load failed: %v\n", err)
		} else {
			fmt.Printf("Load succeeded!\n")
			// Show a success modal
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Successfully loaded: %s", fileName)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(list, true)
				})
			app.SetRoot(modal, true)
		}
	})

	// Add escape to quit
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	})

	// Run the app
	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}
