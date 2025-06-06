package components

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowLibraryContent displays a specific library topic with its content
func ShowLibraryContent(pages *tview.Pages, title, content string, returnPage string) {
	// Create a text view for displaying the content
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(content).
		SetScrollable(true).
		SetBorder(true).
		SetTitle(title)

	// Add keyboard handler for ESC key
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage(returnPage)
			return nil
		}
		return event
	})

	// Add the text view to the pages and show it
	pageName := "library-content"
	pages.AddPage(pageName, textView, true, true)
	pages.SwitchToPage(pageName)
}

// ShowLibraryTopicsList displays a list of available library topics
func ShowLibraryTopicsList(pages *tview.Pages, topics map[string]string, returnPage string) {
	// Format the library topics with colors
	text := "[yellow][::b]Library Topics[-::-]\n\n"

	// Get all topic names and sort them
	topicNames := make([]string, 0, len(topics))
	for name := range topics {
		topicNames = append(topicNames, name)
	}
	sort.Strings(topicNames)

	// Add each topic to the text
	for _, name := range topicNames {
		desc := topics[name]
		text += fmt.Sprintf("[blue]%s[-]: %s\n", name, desc)
	}

	// Create a text view for displaying the topics
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(text).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Library Topics")

	// Add keyboard handler for ESC key
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage(returnPage)
			return nil
		}
		return event
	})

	// Add the text view to the pages and show it
	pageName := "library-topics"
	pages.AddPage(pageName, textView, true, true)
	pages.SwitchToPage(pageName)
}
