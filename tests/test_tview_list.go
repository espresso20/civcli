package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create a simple list
	list := tview.NewList().
		AddItem("test_save.json", "A test save file", '1', nil).
		AddItem("early_game.json", "Early game save", '2', nil).
		AddItem("advanced_civilization.json", "Advanced civilization", '3', nil)

	list.SetBorder(true).SetTitle("Test Save Files")

	// Set up selection handler
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		app.QueueUpdateDraw(func() {
			// Just show a modal instead of doing file operations
			modal := tview.NewModal().
				SetText(fmt.Sprintf("Would load: %s (index %d)", mainText, index)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					app.SetRoot(list, true)
				})
			app.SetRoot(modal, true)
		})
	})

	// Add escape key to quit
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		log.Fatal(err)
	}
}
