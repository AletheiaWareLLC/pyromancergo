package main

import (
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type CacheUI struct {
	widget.BaseWidget
	cache *vm.Cache
}

func NewCacheUI(c *vm.Cache) fyne.CanvasObject {
	ui := &CacheUI{
		cache: c,
	}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *CacheUI) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	sizeLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lineWidthLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lineCountLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	busWidthLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	tagBitsLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	indexBitsLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	offsetBitsLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
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
	addressLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	operationLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lowerAddressLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	lowerOperationLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})
	form := widget.NewForm(
		widget.NewFormItem("Size", sizeLabel),
		widget.NewFormItem("Line Width", lineWidthLabel),
		widget.NewFormItem("Line Count", lineCountLabel),
		widget.NewFormItem("Bus Width", busWidthLabel),
		widget.NewFormItem("Tag Bits", tagBitsLabel),
		widget.NewFormItem("Index Bits", indexBitsLabel),
		widget.NewFormItem("Offset Bits", offsetBitsLabel),
		widget.NewFormItem("Successful", isSuccessfulLabel),
		widget.NewFormItem("Busy", isBusyLabel),
		widget.NewFormItem("Free", isFreeLabel),
		widget.NewFormItem("Address", addressLabel),
		widget.NewFormItem("Operation", operationLabel),
		widget.NewFormItem("Lower Address", lowerAddressLabel),
		widget.NewFormItem("Lower Operation", lowerOperationLabel),
	)
	lines := ui.cache.Lines()
	width := lines[0].Bus.Size()
	columns := 1 + 1 + width/8 // Index + Tag + LineWidth/BytesPerCell
	table := widget.NewTable(
		func() (int, int) {
			return len(lines) + 1 /*Header*/, columns
		},
		func() fyne.CanvasObject {
			return container.NewMax(widget.NewLabelWithStyle(strings.Repeat("0", 18), fyne.TextAlignCenter, fyne.TextStyle{
				Monospace: true,
			}), NewDataCell())
		},
		func(id widget.TableCellID, item fyne.CanvasObject) {
			objects := item.(*fyne.Container).Objects
			label := objects[0].(*widget.Label)
			data := objects[1].(*DataCell)
			col := id.Col
			row := id.Row
			if row == 0 {
				var text string
				// Header
				switch col {
				case 0:
					text = "Index"
				case 1:
					text = "Tag"
				default:
					o := (col - 2) * 8 // BytesPerCell
					text = fmt.Sprintf("Offset 0x%x (%d)", o, o)
				}
				label.SetText(text)
				label.Show()
				data.Hide()
			} else {
				row--
				switch col {
				case 0:
					label.SetText(fmt.Sprintf("0x%x (%d)", row, row))
					label.Show()
					data.Hide()
				case 1:
					label.SetText(fmt.Sprintf("0x%x", lines[row].Tag()))
					label.Show()
					data.Hide()
				default:
					line := lines[row]
					col -= 2
					data.Bus = line
					data.Address = col * 8
					data.Refresh()
					label.Hide()
					data.Show()
				}
			}
		},
	)
	r := &cacheUIRenderer{
		ui:             ui,
		size:           sizeLabel,
		lineWidth:      lineWidthLabel,
		lineCount:      lineCountLabel,
		busWidth:       busWidthLabel,
		tagBits:        tagBitsLabel,
		indexBits:      indexBitsLabel,
		offsetBits:     offsetBitsLabel,
		isSuccessful:   isSuccessfulLabel,
		isBusy:         isBusyLabel,
		isFree:         isFreeLabel,
		address:        addressLabel,
		operation:      operationLabel,
		lowerAddress:   lowerAddressLabel,
		lowerOperation: lowerOperationLabel,
		information:    form,
		table:          table,
		root:           container.NewBorder(nil, nil, form, nil, table),
	}
	r.Refresh()
	return r
}

func (ui *CacheUI) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*cacheUIRenderer)(nil)

type cacheUIRenderer struct {
	ui             *CacheUI
	size           *widget.Label
	lineWidth      *widget.Label
	lineCount      *widget.Label
	busWidth       *widget.Label
	tagBits        *widget.Label
	indexBits      *widget.Label
	offsetBits     *widget.Label
	isSuccessful   *widget.Label
	isBusy         *widget.Label
	isFree         *widget.Label
	address        *widget.Label
	operation      *widget.Label
	lowerAddress   *widget.Label
	lowerOperation *widget.Label
	information    *widget.Form
	table          *widget.Table
	root           *fyne.Container
}

func (r *cacheUIRenderer) Destroy() {}

func (r *cacheUIRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *cacheUIRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *cacheUIRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *cacheUIRenderer) Refresh() {
	r.size.SetText(fmt.Sprintf("%d", r.ui.cache.Size()))
	r.lineWidth.SetText(fmt.Sprintf("%d", r.ui.cache.LineWidth()))
	r.lineCount.SetText(fmt.Sprintf("%d", len(r.ui.cache.Lines())))
	r.busWidth.SetText(fmt.Sprintf("%d", r.ui.cache.Bus().Size()))
	// TODO show bus contents and flags
	r.tagBits.SetText(fmt.Sprintf("%d", r.ui.cache.TagBits()))
	r.indexBits.SetText(fmt.Sprintf("%d", r.ui.cache.IndexBits()))
	r.offsetBits.SetText(fmt.Sprintf("%d", r.ui.cache.OffsetBits()))
	r.address.SetText(fmt.Sprintf("0x%016x", r.ui.cache.Address()))
	r.isSuccessful.SetText(fmt.Sprintf("%t", r.ui.cache.IsSuccessful()))
	r.isBusy.SetText(fmt.Sprintf("%t", r.ui.cache.IsBusy()))
	r.isFree.SetText(fmt.Sprintf("%t", r.ui.cache.IsFree()))
	r.operation.SetText(r.ui.cache.Operation().String())
	r.lowerOperation.SetText(r.ui.cache.LowerOperation().String())
	r.lowerAddress.SetText(fmt.Sprintf("0x%016x", r.ui.cache.LowerAddress()))
	r.information.Refresh()
	r.table.Refresh()
	r.root.Refresh()
}
