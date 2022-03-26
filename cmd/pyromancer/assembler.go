package main

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/assembler"
	"aletheiaware.com/flamego/vm"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"strings"
)

type PyromancerAssembler struct {
	App     fyne.App
	Window  fyne.Window
	Editor  *widget.Entry
	Address *widget.Label
}

func main() {
	pa := &PyromancerAssembler{
		App: app.New(),
		Editor: &widget.Entry{
			MultiLine: true,
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
			Wrapping: fyne.TextWrapBreak,
		},
		Address: &widget.Label{
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
			Wrapping: fyne.TextWrapOff,
		},
	}
	pa.Editor.OnChanged = func(t string) {
		a := assembler.NewAssembler()
		if _, err := a.ReadFrom(strings.NewReader(t)); err != nil {
			if e, ok := err.(*assembler.Error); ok {
				// TODO scroll to and highlight line in editor
				log.Println("Line:", e.Line)
			}
			log.Println(err)
			return
		}

		var buffer bytes.Buffer
		count, err := a.WriteTo(&buffer)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(count, "bytes")
		// TODO pa.UpdateBinary(a) -> pa.Binary.SetText(buffer.String())
		pa.UpdateAddress(a)
	}
	pa.Window = pa.App.NewWindow("Pyromancer Assembler")
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			dialog.ShowError(err, pa.Window)
			return
		}
		pa.Open(f.Name(), f)
	}
	pa.Window.SetContent(container.NewBorder(widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			pa.Editor.SetText("")
			pa.Window.SetTitle("Pyromancer Assembler")
		}),
		widget.NewToolbarAction(theme.FileIcon(), func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, pa.Window)
					return
				}
				if reader != nil {
					pa.Open(reader.URI().Name(), reader)
				}
			}, pa.Window)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".fas"}))
			fd.Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, pa.Window)
					return
				}
				if writer == nil {
					return
				}
				if _, err := writer.Write([]byte(pa.Editor.Text)); err != nil {
					dialog.ShowError(err, pa.Window)
					return
				}
			}, pa.Window)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".fas"}))
			fd.Show()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			m, err := pa.NewMachine()
			if err != nil {
				dialog.ShowError(err, pa.Window)
				return
			}
			// Signal the first context of the first core
			m.Processor.Signal(0)
			NewInspector(pa.App, m)
		}),
		widget.NewToolbarAction(theme.ComputerIcon(), func() {
			m, err := pa.NewMachine()
			if err != nil {
				dialog.ShowError(err, pa.Window)
				return
			}
			// Signal the first context of the first core
			m.Processor.Signal(0)
			pe := NewEmulator(pa.App, m)
			go pe.ClockContinuously()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			log.Println("Not Yet Implemented")
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Not Yet Implemented")
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			log.Println("Not Yet Implemented")
		}),
	), nil, nil, nil, container.NewHSplit(pa.Editor, container.NewScroll(pa.Address))))
	pa.Window.CenterOnScreen()
	pa.Window.Resize(fyne.NewSize(800, 600))
	pa.Window.ShowAndRun()
}

func (pa *PyromancerAssembler) Open(name string, reader io.ReadCloser) {
	pa.Window.SetTitle("Pyromancer Assembler - " + name)
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		dialog.ShowError(err, pa.Window)
		return
	}
	pa.Editor.SetText(string(data))
}

func (pa *PyromancerAssembler) NewMachine() (*vm.Machine, error) {
	a := assembler.NewAssembler()
	if _, err := a.ReadFrom(strings.NewReader(pa.Editor.Text)); err != nil {
		if e, ok := err.(*assembler.Error); ok {
			// TODO scroll to and highlight line in editor
			log.Println("Line:", e.Line)
		}
		return nil, err
	}

	m := vm.NewMachine()

	// Add storage device
	m.Processor.AddDevice(vm.NewFileStorage(m.Memory, flamego.DeviceControlBlockAddress)) // TODO processor can set control block offset

	// Add display device
	m.Processor.AddDevice(vm.NewDisplay(m.Memory, flamego.DeviceControlBlockAddress+flamego.DeviceControlBlockSize, 320, 240)) // TODO processor can set control block offset

	var buffer bytes.Buffer
	count, err := a.WriteTo(&buffer)
	if err != nil {
		return nil, err
	}
	log.Println(count, "bytes")
	copy(m.Memory.Data(), buffer.Bytes())

	pa.UpdateAddress(a)

	return m, nil
}

func (pa *PyromancerAssembler) UpdateAddress(a assembler.Assembler) {
	var addresses string
	for _, a := range a.Addressables() {
		addresses += fmt.Sprintf("0x%016x", a.AbsoluteAddress())
		if s, ok := a.(fmt.Stringer); ok {
			addresses += ": "
			addresses += s.String()
		}
		addresses += "\n"
	}
	pa.Address.SetText(addresses)
}
