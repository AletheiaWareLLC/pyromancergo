package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	t := &canvas.Text{
		Alignment: fyne.TextAlignCenter,
		Color:     theme.TextColor(),
		TextSize:  theme.TextSize(),
		TextStyle: fyne.TextStyle{
			Monospace: true,
		},
	}
	r := &buttonRenderer{
		b:    b,
		text: t,
		root: container.NewPadded(t),
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
	b    *Button
	text *canvas.Text
	root *fyne.Container
}

func (r *buttonRenderer) Destroy() {}

func (r *buttonRenderer) Layout(size fyne.Size) {
	r.root.Resize(size)
}

func (r *buttonRenderer) MinSize() fyne.Size {
	return r.root.MinSize()
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.root}
}

func (r *buttonRenderer) Refresh() {
	r.text.Text = r.b.Name
	r.text.Refresh()
}
