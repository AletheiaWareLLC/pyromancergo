package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CoreUI struct {
	widget.BaseWidget
	core *vm.Core
}

func NewCoreUI(c *vm.Core) fyne.CanvasObject {
	ui := &CoreUI{
		core: c,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *CoreUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	nextLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lockHolderLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	requiresLockLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	acquiredLockLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	loadRegister0Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	loadRegister1Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	loadRegister2Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	loadRegister3Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	executeRegister0Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	executeRegister1Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	formatRegister0Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	formatRegister1Label := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("Next Context", nextLabel),
		widget.NewFormItem("Lock Holder", lockHolderLabel),
		widget.NewFormItem("Requires Lock", requiresLockLabel),
		widget.NewFormItem("Acquired Lock", acquiredLockLabel),
		widget.NewFormItem("Load Register 0", loadRegister0Label),
		widget.NewFormItem("Load Register 1", loadRegister1Label),
		widget.NewFormItem("Load Register 2", loadRegister2Label),
		widget.NewFormItem("Load Register 3", loadRegister3Label),
		widget.NewFormItem("Execute Register 0", executeRegister0Label),
		widget.NewFormItem("Execute Register 1", executeRegister1Label),
		widget.NewFormItem("Format Register 0", formatRegister0Label),
		widget.NewFormItem("Format Register 1", formatRegister1Label),
	)
	r := &coreUIRenderer{
		ui:               ui,
		next:             nextLabel,
		lockHolder:       lockHolderLabel,
		requiresLock:     requiresLockLabel,
		acquiredLock:     acquiredLockLabel,
		loadRegister0:    loadRegister0Label,
		loadRegister1:    loadRegister1Label,
		loadRegister2:    loadRegister2Label,
		loadRegister3:    loadRegister3Label,
		executeRegister0: executeRegister0Label,
		executeRegister1: executeRegister1Label,
		formatRegister0:  formatRegister0Label,
		formatRegister1:  formatRegister1Label,
		information:      form,
		root:             container.NewMax(form),
	}
	r.Refresh()
	return r
}

func (ui *CoreUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*coreUIRenderer)(nil)

type coreUIRenderer struct {
	ui               *CoreUI
	next             *widget.Label
	lockHolder       *widget.Label
	requiresLock     *widget.Label
	acquiredLock     *widget.Label
	loadRegister0    *widget.Label
	loadRegister1    *widget.Label
	loadRegister2    *widget.Label
	loadRegister3    *widget.Label
	executeRegister0 *widget.Label
	executeRegister1 *widget.Label
	formatRegister0  *widget.Label
	formatRegister1  *widget.Label
	information      *widget.Form
	root             *fyne.Container
}

func (r *coreUIRenderer) Destroy() {}

func (r *coreUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *coreUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *coreUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *coreUIRenderer) Refresh() {
	r.next.SetText(fmt.Sprintf("%d", r.ui.core.NextContext()))
	holder := "-"
	if h := r.ui.core.LockHolder(); h >= 0 {
		holder = fmt.Sprintf("%d", h)
	}
	r.lockHolder.SetText(holder)
	r.requiresLock.SetText(fmt.Sprintf("%t", r.ui.core.RequiresLock()))
	r.acquiredLock.SetText(fmt.Sprintf("%t", r.ui.core.AcquiredLock()))
	r.loadRegister0.SetText(fmt.Sprintf("0x%016x", r.ui.core.LoadRegister0()))
	r.loadRegister1.SetText(fmt.Sprintf("0x%016x", r.ui.core.LoadRegister1()))
	r.loadRegister2.SetText(fmt.Sprintf("0x%016x", r.ui.core.LoadRegister2()))
	r.loadRegister3.SetText(fmt.Sprintf("0x%016x", r.ui.core.LoadRegister3()))
	r.executeRegister0.SetText(fmt.Sprintf("0x%016x", r.ui.core.ExecuteRegister0()))
	r.executeRegister1.SetText(fmt.Sprintf("0x%016x", r.ui.core.ExecuteRegister1()))
	r.formatRegister0.SetText(fmt.Sprintf("0x%016x", r.ui.core.FormatRegister0()))
	r.formatRegister1.SetText(fmt.Sprintf("0x%016x", r.ui.core.FormatRegister1()))
	r.information.Refresh()
	r.root.Refresh()
}
