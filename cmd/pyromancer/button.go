package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var _ (fyne.Tappable) = (*Button)(nil)

type Button struct {
	widget.BaseWidget
	Name     string
	OnTapped func()
}

func NewButton(name string, onTapped func()) fyne.CanvasObject {
	b := &Button{
		Name:     name,
		OnTapped: onTapped,
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *Button) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	r := &buttonRenderer{
		b: b,
		label: widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
	}
	r.Refresh()
	return r
}

func (b *Button) MinSize() fyne.Size {
	b.ExtendBaseWidget(b)
	return b.BaseWidget.MinSize()
}

func (b *Button) Tapped(*fyne.PointEvent) {
	if t := b.OnTapped; t != nil {
		t()
	}
}

var _ fyne.WidgetRenderer = (*buttonRenderer)(nil)

type buttonRenderer struct {
	b     *Button
	label *widget.Label
}

func (r *buttonRenderer) Destroy() {}

func (r *buttonRenderer) Layout(size fyne.Size) {
	r.label.Resize(size)
}

func (r *buttonRenderer) MinSize() fyne.Size {
	return r.label.MinSize()
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.label}
}

func (r *buttonRenderer) Refresh() {
	r.label.SetText(r.b.Name)
}
