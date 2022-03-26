package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProcessorUI struct {
	widget.BaseWidget
	processor *vm.Processor
}

func NewProcessorUI(p *vm.Processor) fyne.CanvasObject {
	ui := &ProcessorUI{
		processor: p,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *ProcessorUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	stateLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lockHolderLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("State", stateLabel),
		widget.NewFormItem("Lock Holder", lockHolderLabel),
	)
	r := &processorUIRenderer{
		ui:          ui,
		state:       stateLabel,
		lockHolder:  lockHolderLabel,
		information: form,
		root:        container.NewMax(form),
	}
	r.Refresh()
	return r
}

func (ui *ProcessorUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*processorUIRenderer)(nil)

type processorUIRenderer struct {
	ui          *ProcessorUI
	state       *widget.Label
	lockHolder  *widget.Label
	information *widget.Form
	root        *fyne.Container
}

func (r *processorUIRenderer) Destroy() {}

func (r *processorUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *processorUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *processorUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *processorUIRenderer) Refresh() {
	state := "running"
	if r.ui.processor.HasHalted() {
		state = "halted"
	}
	r.state.SetText(state)
	holder := "-"
	if h := r.ui.processor.LockHolder(); h >= 0 {
		holder = fmt.Sprintf("%d", h)
	}
	r.lockHolder.SetText(holder)
	r.information.Refresh()
	r.root.Refresh()
}
