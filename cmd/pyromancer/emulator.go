package main

import (
	"aletheiaware.com/flamego/vm"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type PyromancerEmulator struct {
	Window  fyne.Window
	Machine *vm.Machine
	Raster  *canvas.Raster
	stop    chan struct{}
}

func NewEmulator(a fyne.App, m *vm.Machine) *PyromancerEmulator {
	pe := &PyromancerEmulator{
		Window:  a.NewWindow("Pyromancer Emulator"),
		Machine: m,
	}
	// TODO fileEntry := NewFileEntry(pe.Machine.Processor.Device(0).(*vm.FileStorage), pe.Window)
	pe.Raster = NewDisplayRaster(pe.Machine.Processor.Device(1).(*vm.Display))
	pe.Raster.SetMinSize(fyne.NewSize(320, 240))

	pe.Window.SetContent(container.NewCenter(pe.Raster))
	pe.Window.SetOnClosed(pe.Stop)
	pe.Window.CenterOnScreen()
	pe.Window.Resize(fyne.NewSize(800, 600))
	pe.Window.Show()
	return pe
}

func (pe *PyromancerEmulator) Clock() {
	pe.Machine.Clock()

	// Update UI
	pe.Raster.Refresh()
}

func (pe *PyromancerEmulator) ClockContinuously() {
	pe.Stop()
	pe.stop = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-pe.stop:
				close(pe.stop)
				pe.stop = nil
				return
			default:
				pe.Clock()
				if pe.Machine.Processor.HasHalted() {
					pe.Stop()
				}
			}
		}
	}()
}

func (pe *PyromancerEmulator) Stop() {
	if s := pe.stop; s != nil {
		select {
		case s <- struct{}{}:
		}
	}
}
