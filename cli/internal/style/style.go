package style

import "github.com/charmbracelet/lipgloss"

var (
	// Description is the style for user descriptions of items.
	Description = lipgloss.NewStyle().
			Faint(true)

	// Instructions is the style for user interface instructions.
	Instructions = Description

	// Selected is the style for something that is currently selected in some
	// way, such as a list item.
	Selected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0080ff"))

	// Unselected is the style for something that is selectable but is not
	// currently selected.
	Unselected = Selected.
			Copy().
			Faint(true)

	// Spinner is the style to apply to spinner components.
	Spinner = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fcba03"))
)
