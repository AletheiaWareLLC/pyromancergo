package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DisplayUI struct {
	widget.BaseWidget
	display *vm.Display
	window  fyne.Window
}

func NewDisplayUI(s *vm.Display, w fyne.Window) fyne.CanvasObject {
	ui := &DisplayUI{
		display: s,
		window:  w,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *DisplayUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	memoryOffsetLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	memoryOperationLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	isBusyLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	operationLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	commandLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	controllerLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	parameterLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	deviceAddressLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	memoryAddressLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("Memory Offset", memoryOffsetLabel),
		widget.NewFormItem("Memory Operation", memoryOperationLabel),
		widget.NewFormItem("Busy", isBusyLabel),
		widget.NewFormItem("Operation", operationLabel),
		widget.NewFormItem("Command", commandLabel),
		widget.NewFormItem("Controller", controllerLabel),
		widget.NewFormItem("Parameter", parameterLabel),
		widget.NewFormItem("Device Address", deviceAddressLabel),
		widget.NewFormItem("Memory Address", memoryAddressLabel),
	)
	raster := NewDisplayRaster(ui.display)
	raster.SetMinSize(fyne.NewSize(320, 240))
	r := &displayUIRenderer{
		ui:              ui,
		memoryOffset:    memoryOffsetLabel,
		memoryOperation: memoryOperationLabel,
		isBusy:          isBusyLabel,
		operation:       operationLabel,
		command:         commandLabel,
		controller:      controllerLabel,
		parameter:       parameterLabel,
		deviceAddress:   deviceAddressLabel,
		memoryAddress:   memoryAddressLabel,
		information:     form,
		raster:          raster,
		root:            container.NewBorder(form, nil, nil, nil, container.NewCenter(raster)),
	}
	r.Refresh()
	return r
}

func (ui *DisplayUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*displayUIRenderer)(nil)

type displayUIRenderer struct {
	ui              *DisplayUI
	memoryOffset    *widget.Label
	memoryOperation *widget.Label
	isBusy          *widget.Label
	operation       *widget.Label
	command         *widget.Label
	controller      *widget.Label
	parameter       *widget.Label
	deviceAddress   *widget.Label
	memoryAddress   *widget.Label
	information     *widget.Form
	raster          *canvas.Raster
	root            *fyne.Container
}

func (r *displayUIRenderer) Destroy() {}

func (r *displayUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *displayUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *displayUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *displayUIRenderer) Refresh() {
	r.memoryOffset.SetText(fmt.Sprintf("0x%016x", r.ui.display.MemoryOffset()))
	var op string
	o := r.ui.display.MemoryOperation()
	switch o {
	case vm.ReadCommand:
		op = "Read Command"
	case vm.ReadDeviceAddress:
		op = "Read Device Address"
	case vm.ReadMemoryAddress:
		op = "Read Memory Address"
	default:
		op = o.String()
	}
	r.memoryOperation.SetText(op)
	r.isBusy.SetText(fmt.Sprintf("%t", r.ui.display.IsBusy()))
	r.operation.SetText(r.ui.display.Operation().String())
	r.command.SetText(fmt.Sprintf("0x%016x", r.ui.display.Command()))
	r.controller.SetText(fmt.Sprintf("0x%02x", r.ui.display.Controller()))
	r.parameter.SetText(fmt.Sprintf("0x%012x", r.ui.display.Parameter()))
	r.deviceAddress.SetText(fmt.Sprintf("0x%016x", r.ui.display.DeviceAddress()))
	r.memoryAddress.SetText(fmt.Sprintf("0x%016x", r.ui.display.MemoryAddress()))
	r.information.Refresh()
	r.raster.Refresh()
	r.root.Refresh()
}
