package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ToolbarCycleCount struct {
	entry *widget.Entry
}

func NewToolbarCycleCount(entry *widget.Entry) widget.ToolbarItem {
	return &ToolbarCycleCount{
		entry: entry,
	}
}

func (t *ToolbarCycleCount) ToolbarObject() fyne.CanvasObject {
	return container.NewPadded(container.NewHBox(
		widget.NewLabelWithStyle("Cycle", fyne.TextAlignTrailing, fyne.TextStyle{
			Bold:      true,
			Monospace: true,
		}),
		t.entry,
	))
}
