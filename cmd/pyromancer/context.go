package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ContextUI struct {
	widget.BaseWidget
	context *vm.Context
}

func NewContextUI(c *vm.Context) fyne.CanvasObject {
	ui := &ContextUI{
		context: c,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *ContextUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	validLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	asleepLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	sleepCyclesLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	interruptedLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	nextInterruptLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	signalledLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	retryingLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	alignedLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	statusLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	opcodeLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	instructionLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	requiresLockLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	acquiredLockLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("Valid", validLabel),
		widget.NewFormItem("Asleep", asleepLabel),
		widget.NewFormItem("Interrupted", interruptedLabel),
		widget.NewFormItem("Signalled", signalledLabel),
		widget.NewFormItem("Retrying", retryingLabel),
		widget.NewFormItem("Aligned", alignedLabel),
		widget.NewFormItem("Requires Lock", requiresLockLabel),
		widget.NewFormItem("Acquired Lock", acquiredLockLabel),
		widget.NewFormItem("Status", statusLabel),
		widget.NewFormItem("Cycles Asleep", sleepCyclesLabel),
		widget.NewFormItem("Next Interrupt", nextInterruptLabel),
		widget.NewFormItem("Opcode", opcodeLabel),
		widget.NewFormItem("Instruction", instructionLabel),
	)
	registers := NewRegisterUI(ui.context.ReadRegister)
	r := &contextUIRenderer{
		ui:            ui,
		valid:         validLabel,
		asleep:        asleepLabel,
		interrupted:   interruptedLabel,
		signalled:     signalledLabel,
		retrying:      retryingLabel,
		aligned:       alignedLabel,
		requiresLock:  requiresLockLabel,
		acquiredLock:  acquiredLockLabel,
		status:        statusLabel,
		sleepCycles:   sleepCyclesLabel,
		nextInterrupt: nextInterruptLabel,
		opcode:        opcodeLabel,
		instruction:   instructionLabel,
		information:   form,
		registers:     registers,
		root:          container.NewHSplit(form, registers),
	}
	r.Refresh()
	return r
}

func (ui *ContextUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*contextUIRenderer)(nil)

type contextUIRenderer struct {
	ui            *ContextUI
	valid         *widget.Label
	asleep        *widget.Label
	interrupted   *widget.Label
	signalled     *widget.Label
	retrying      *widget.Label
	aligned       *widget.Label
	requiresLock  *widget.Label
	acquiredLock  *widget.Label
	status        *widget.Label
	sleepCycles   *widget.Label
	nextInterrupt *widget.Label
	opcode        *widget.Label
	instruction   *widget.Label
	information   *widget.Form
	registers     fyne.CanvasObject
	root          *container.Split
}

func (r *contextUIRenderer) Destroy() {}

func (r *contextUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *contextUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *contextUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *contextUIRenderer) Refresh() {
	r.valid.SetText(fmt.Sprintf("%t", r.ui.context.IsValid()))
	r.asleep.SetText(fmt.Sprintf("%t", r.ui.context.IsAsleep()))
	r.interrupted.SetText(fmt.Sprintf("%t", r.ui.context.IsInterrupted()))
	r.signalled.SetText(fmt.Sprintf("%t", r.ui.context.IsSignalled()))
	r.retrying.SetText(fmt.Sprintf("%t", r.ui.context.IsRetrying()))
	r.aligned.SetText(fmt.Sprintf("%t", r.ui.context.IsAligned()))
	r.requiresLock.SetText(fmt.Sprintf("%t", r.ui.context.RequiresLock()))
	r.acquiredLock.SetText(fmt.Sprintf("%t", r.ui.context.AcquiredLock()))
	r.status.SetText(r.ui.context.Status())
	r.sleepCycles.SetText(fmt.Sprintf("%d", r.ui.context.SleepCycles()))
	inter := "-"
	if i := r.ui.context.NextInterrupt(); i >= 0 {
		inter = fmt.Sprintf("0x%04x", uint16(i))
	}
	r.nextInterrupt.SetText(inter)
	opcode := "-"
	if o := r.ui.context.Opcode(); o != 0 {
		opcode = fmt.Sprintf("0x%08x", o)
	}
	r.opcode.SetText(opcode)
	r.instruction.SetText(r.ui.context.InstructionString())
	r.information.Refresh()
	r.registers.Refresh()
	r.root.Refresh()
}
