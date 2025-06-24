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
	fmt.Println("Testing load game with automatic Enter simulation...")

	app := tview.NewApplication()
	engine := game.NewGameEngine(&DummyDisplay{})

	saveDir := "./data/saves"
	files, err := os.ReadDir(saveDir)
	if err != nil {
		panic(err)
	}

	list := tview.NewList()
	list.SetBorder(true).SetTitle("Auto Load Test")

	var saveFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			saveFiles = append(saveFiles, file.Name())
			list.AddItem(file.Name(), fmt.Sprintf("Size: %d", 100), 0, nil)
		}
	}

	fmt.Printf("Found %d save files\n", len(saveFiles))

	// Simulated Enter key press
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		fmt.Printf("*** ENTER PRESSED: index=%d, file=%s ***\n", index, mainText)

		if index >= len(saveFiles) {
			fmt.Printf("ERROR: Invalid index\n")
			return
		}

		fileName := strings.TrimSuffix(saveFiles[index], ".json")
		fmt.Printf("Loading: %s\n", fileName)

		// The critical section: does LoadGame hang/crash?
		fmt.Printf("Before LoadGame call...\n")
		err := engine.LoadGame(fileName)
		fmt.Printf("After LoadGame call, err=%v\n", err)

		if err != nil {
			fmt.Printf("Load failed: %v\n", err)
		} else {
			fmt.Printf("Load succeeded!\n")
		}

		fmt.Printf("Stopping app...\n")
		app.Stop()
	})

	// Auto-trigger after 1 second
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("Auto-triggering Enter on first save file...\n")
		app.QueueUpdateDraw(func() {
			if len(saveFiles) > 0 {
				// Directly call the selected func to simulate Enter key
				list.SetCurrentItem(0)
			}
		})

		// If the app doesn't respond in 5 seconds, force quit
		time.Sleep(5 * time.Second)
		fmt.Printf("Timeout reached - force stopping app\n")
		app.Stop()
	}()

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	})

	fmt.Printf("Starting app...\n")
	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
	fmt.Printf("App finished.\n")
}
