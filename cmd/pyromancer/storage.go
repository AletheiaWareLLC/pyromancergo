package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FileStorageUI struct {
	widget.BaseWidget
	storage *vm.FileStorage
	window  fyne.Window
}

func NewFileStorageUI(s *vm.FileStorage, w fyne.Window) fyne.CanvasObject {
	ui := &FileStorageUI{
		storage: s,
		window:  w,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *FileStorageUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	fileEntry := NewFileEntry(ui.storage, ui.window)
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
		widget.NewFormItem("File", fileEntry),
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
	r := &fileStorageUIRenderer{
		ui:              ui,
		file:            fileEntry,
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
		root:            container.NewMax(form),
	}
	r.Refresh()
	return r
}

func (ui *FileStorageUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*fileStorageUIRenderer)(nil)

type fileStorageUIRenderer struct {
	ui              *FileStorageUI
	file            *widget.Entry
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
	root            *fyne.Container
}

func (r *fileStorageUIRenderer) Destroy() {}

func (r *fileStorageUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *fileStorageUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *fileStorageUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *fileStorageUIRenderer) Refresh() {
	var n string
	if f := r.ui.storage.File(); f != nil {
		n = f.Name()
	}
	r.file.SetText(n)
	r.memoryOffset.SetText(fmt.Sprintf("0x%016x", r.ui.storage.MemoryOffset()))
	var op string
	o := r.ui.storage.MemoryOperation()
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
	r.isBusy.SetText(fmt.Sprintf("%t", r.ui.storage.IsBusy()))
	r.operation.SetText(r.ui.storage.Operation().String())
	r.command.SetText(fmt.Sprintf("0x%016x", r.ui.storage.Command()))
	r.controller.SetText(fmt.Sprintf("0x%02x", r.ui.storage.Controller()))
	r.parameter.SetText(fmt.Sprintf("0x%012x", r.ui.storage.Parameter()))
	r.deviceAddress.SetText(fmt.Sprintf("0x%016x", r.ui.storage.DeviceAddress()))
	r.memoryAddress.SetText(fmt.Sprintf("0x%016x", r.ui.storage.MemoryAddress()))
	r.information.Refresh()
	r.root.Refresh()
}
