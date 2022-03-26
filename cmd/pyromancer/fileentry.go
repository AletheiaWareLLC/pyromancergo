package main

import (
	"aletheiaware.com/flamego/vm"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewFileEntry(s *vm.FileStorage, w fyne.Window) *widget.Entry {
	fileEntry := widget.NewEntry()
	//fileEntry.TextAlign = fyne.TextAlignCenter
	fileEntry.TextStyle = fyne.TextStyle{
		Monospace: true,
	}
	fileEntry.OnSubmitted = func(t string) {
		if err := s.Open(t); err != nil {
			dialog.ShowError(err, w)
		}
	}
	b := widget.NewButtonWithIcon("", theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader != nil {
				if err := reader.Close(); err != nil {
					dialog.ShowError(err, w)
					return
				}
				if err := s.Open(reader.URI().Path()); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bin"}))
		fd.Show()
	})
	b.Importance = widget.LowImportance
	fileEntry.ActionItem = b
	return fileEntry
}
