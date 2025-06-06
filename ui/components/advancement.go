package components

import (
	"fmt"

	"github.com/rivo/tview"
)

// ShowAgeAdvancement displays a modal with animation to celebrate age advancement
func ShowAgeAdvancement(pages *tview.Pages, newAge string, returnPage string) {
	message := fmt.Sprintf("[yellow]⭐ ADVANCEMENT ⭐[-]\n[green]Your civilization has advanced to the [::b]%s[-::-]", newAge)

	// Create a modal to show the advancement
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Continue"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.SwitchToPage(returnPage)
		})

	// Add the modal to the pages and show it
	pageName := "advancement"
	pages.AddPage(pageName, modal, true, true)
	pages.SwitchToPage(pageName)
}
