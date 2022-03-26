package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type MemoryUI struct {
	widget.BaseWidget
	memory *vm.Memory
}

func NewMemoryUI(c *vm.Memory) fyne.CanvasObject {
	ui := &MemoryUI{
		memory: c,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *MemoryUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	sizeLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	busWidthLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	addressLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	isSuccessfulLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	isBusyLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	isFreeLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	operationLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("Size", sizeLabel),
		widget.NewFormItem("Bus Width", busWidthLabel),
		widget.NewFormItem("Address", addressLabel),
		widget.NewFormItem("Successful", isSuccessfulLabel),
		widget.NewFormItem("Busy", isBusyLabel),
		widget.NewFormItem("Free", isFreeLabel),
		widget.NewFormItem("Operation", operationLabel),
	)
	data := ui.memory.Data()
	table := widget.NewTable(
		func() (int, int) {
			return int(ui.memory.Size())/8 + 1 /*Header*/, 2
		},
		func() fyne.CanvasObject {
			// 0x0000000000000000
			return widget.NewLabelWithStyle(strings.Repeat("0", 18), fyne.TextAlignCenter, fyne.TextStyle{
				Monospace: true,
			})
		},
		func(id widget.TableCellID, item fyne.CanvasObject) {
			col := id.Col
			row := id.Row
			var text string
			if row == 0 {
				//Header
				switch col {
				case 0:
					text = "Address"
				default:
					text = "Data"
				}
			} else {
				row--
				address := uint64(row) * 8
				text += "0x"
				switch col {
				case 0:
					// Address
					text += fmt.Sprintf("%x (%d)", address, address)
				case 1:
					// Value
					for i := uint64(0); i < 8; i += 2 {
						text += fmt.Sprintf("%02x%02x", data[address+i], data[address+i+1])
					}
				}
			}
			item.(*widget.Label).SetText(text)
		},
	)
	r := &memoryUIRenderer{
		ui:           ui,
		size:         sizeLabel,
		busWidth:     busWidthLabel,
		address:      addressLabel,
		isSuccessful: isSuccessfulLabel,
		isBusy:       isBusyLabel,
		isFree:       isFreeLabel,
		operation:    operationLabel,
		information:  form,
		table:        table,
		root:         container.NewBorder(nil, nil, form, nil, table),
	}
	r.Refresh()
	return r
}

func (ui *MemoryUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*memoryUIRenderer)(nil)

type memoryUIRenderer struct {
	ui           *MemoryUI
	size         *widget.Label
	busWidth     *widget.Label
	address      *widget.Label
	isSuccessful *widget.Label
	isBusy       *widget.Label
	isFree       *widget.Label
	operation    *widget.Label
	information  *widget.Form
	table        *widget.Table
	root         *fyne.Container
}

func (r *memoryUIRenderer) Destroy() {}

func (r *memoryUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *memoryUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *memoryUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *memoryUIRenderer) Refresh() {
	r.size.SetText(fmt.Sprintf("%d", r.ui.memory.Size()))
	r.busWidth.SetText(fmt.Sprintf("%d", r.ui.memory.Bus().Size()))
	// TODO show bus contents and flags
	r.address.SetText(fmt.Sprintf("0x%016x", r.ui.memory.Address()))
	r.isSuccessful.SetText(fmt.Sprintf("%t", r.ui.memory.IsSuccessful()))
	r.isBusy.SetText(fmt.Sprintf("%t", r.ui.memory.IsBusy()))
	r.isFree.SetText(fmt.Sprintf("%t", r.ui.memory.IsFree()))
	r.operation.SetText(r.ui.memory.Operation().String())
	r.information.Refresh()
	r.table.Refresh()
	r.root.Refresh()
}
