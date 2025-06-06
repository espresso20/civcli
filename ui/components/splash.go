package components

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	devName = "Espresso"
	version = "v1.0.0"
	wikiURL = "https://github.com/espresso20/civcli/wiki"
)

// openURL opens the provided URL in the default browser
func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

// ShowSplashScreen creates and runs the splash screen
func ShowSplashScreen() error {
	app := tview.NewApplication()

	banner := `
   ____ _              ___    _ _      
  / __|| | __ _ _     |_ _|__| | | ___ 
 | |   | | \ \ / /     |  / _  | |/ _ \
 | |___| |  \ V /      | | (_| | |  __/
  \____|_|   \_/      |___\__,_|_|\___|
`

	intro := "Welcome to Civ-Idle!\n\n" +
		"Civ-Idle is a text-based civilization builder.\n" +
		"Grow your tribe, manage resources, build, and research through the ages.\n\n" +
		"Check out the GitHub wiki for information on how to play:\n"

	// Create a clickable link for the wiki
	wikiLink := fmt.Sprintf("[blue::u]%s[-:-:-]", wikiURL)

	text := banner + "\n" + intro + "\n" + wikiLink

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(text)

	// Make the text view capture mouse events
	textView.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick {
			_, y := event.Position()
			if y == strings.Count(banner+"\n"+intro+"\n", "\n") {
				openURL(wikiURL)
				return action, nil
			}
		}
		return action, event
	})

	// Create button layout
	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	// Create Start button
	startButton := tview.NewButton("Start").
		SetSelectedFunc(func() {
			app.Stop()
		})

	// Create Exit button
	exitButton := tview.NewButton("Exit").
		SetSelectedFunc(func() {
			app.Stop()
			os.Exit(0)
		})

	// Add spacing around buttons
	buttons.AddItem(nil, 0, 1, false)
	buttons.AddItem(startButton, 10, 0, true)
	buttons.AddItem(nil, 2, 0, false)
	buttons.AddItem(exitButton, 10, 0, false)
	buttons.AddItem(nil, 0, 1, false)

	// Create main layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(textView, 0, 1, false).
		AddItem(buttons, 3, 0, true).
		AddItem(nil, 1, 0, false)

	frame := tview.NewFrame(flex).
		AddText("Use Tab/Shift+Tab to navigate, Enter to select", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(devName, false, tview.AlignLeft, tcell.ColorGray).
		AddText(version, false, tview.AlignRight, tcell.ColorGray)

	// Set up keyboard navigation
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			if app.GetFocus() == startButton {
				app.SetFocus(exitButton)
				return nil
			} else if app.GetFocus() == exitButton {
				app.SetFocus(startButton)
				return nil
			}
		case tcell.KeyBacktab:
			if app.GetFocus() == startButton {
				app.SetFocus(exitButton)
				return nil
			} else if app.GetFocus() == exitButton {
				app.SetFocus(startButton)
				return nil
			}
		case tcell.KeyEnter:
			if app.GetFocus() == startButton {
				app.Stop()
				return nil
			} else if app.GetFocus() == exitButton {
				app.Stop()
				os.Exit(0)
				return nil
			}
		}
		return event
	})

	// Initial focus on Start button
	app.SetFocus(startButton)

	return app.SetRoot(frame, true).Run()
}
