package components

import (
	"github.com/rivo/tview"
)

// FormatMessage formats a message based on the specified style
func FormatMessage(message string, style string) string {
	// Format the message based on style
	formattedMessage := message
	switch style {
	case "error":
		formattedMessage = "[red]ERROR: " + message + "[-]"
	case "success":
		formattedMessage = "[green]SUCCESS: " + message + "[-]"
	case "info":
		formattedMessage = "[blue]INFO: " + message + "[-]"
	case "warning":
		formattedMessage = "[yellow]WARNING: " + message + "[-]"
	}

	return formattedMessage
}

// AppendToOutput appends a message to the output text view and scrolls to the end
func AppendToOutput(output *tview.TextView, message string) {
	// Get current text and append
	currentText := output.GetText(false)
	if currentText != "" {
		currentText += "\n"
	}
	output.SetText(currentText + message)

	// Scroll to the bottom
	output.ScrollToEnd()
}
