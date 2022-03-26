package main

import (
	"aletheiaware.com/flamego"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

var (
	red    = color.NRGBA{0xff, 0, 0, 0xff}
	yellow = color.NRGBA{0xff, 0xff, 0, 0xff}
	green  = color.NRGBA{0, 0xff, 0, 0xff}
)

type DataCell struct {
	widget.BaseWidget
	Bus     flamego.Bus
	Address int
}

func NewDataCell() fyne.CanvasObject {
	ui := &DataCell{}
	ui.ExtendBaseWidget(ui)
	return ui
}

func (ui *DataCell) CreateRenderer() fyne.WidgetRenderer {
	ui.ExtendBaseWidget(ui)
	bytes := make([]fyne.CanvasObject, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = &canvas.Text{
			Alignment: fyne.TextAlignCenter,
			Color:     theme.TextColor(),
			Text:      "00",
			TextSize:  theme.TextSize(),
			TextStyle: fyne.TextStyle{
				Monospace: true,
			},
		}
	}
	r := &dataCellRenderer{
		ui:    ui,
		bytes: bytes,
	}
	r.Refresh()
	return r
}

func (ui *DataCell) MinSize() fyne.Size {
	ui.ExtendBaseWidget(ui)
	return ui.BaseWidget.MinSize()
}

var _ fyne.WidgetRenderer = (*dataCellRenderer)(nil)

type dataCellRenderer struct {
	ui    *DataCell
	bytes []fyne.CanvasObject
}

func (r *dataCellRenderer) Destroy() {}

func (r *dataCellRenderer) Layout(size fyne.Size) {
	w := size.Width
	h := size.Height
	cw := w / 8
	for i := 0; i < 8; i++ {
		r.bytes[i].Move(fyne.NewPos(float32(i)*cw, 0))
		r.bytes[i].Resize(fyne.NewSize(cw, h))
	}
}

func (r *dataCellRenderer) MinSize() fyne.Size {
	var s fyne.Size
	for i := 0; i < 8; i++ {
		s.Add(r.bytes[i].MinSize())
	}
	return s
}

func (r *dataCellRenderer) Objects() []fyne.CanvasObject {
	return r.bytes
}

func (r *dataCellRenderer) Refresh() {
	if r.ui.Bus == nil {
		return
	}
	for i := 0; i < 8; i++ {
		index := r.ui.Address + i
		cell := r.bytes[i].(*canvas.Text)
		if !r.ui.Bus.IsValid(index) {
			cell.Color = red
		} else if r.ui.Bus.IsDirty(index) {
			cell.Color = yellow
		} else {
			cell.Color = green
		}
		cell.Text = fmt.Sprintf("%02x", r.ui.Bus.Read(index))
		cell.Refresh()
	}
}
