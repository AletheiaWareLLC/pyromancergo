package main

import (
	"aletheiaware.com/flamego/vm"
	"fyne.io/fyne/v2/canvas"
	"image"
	"image/color"
	"image/draw"
)

func NewDisplayRaster(d *vm.Display) *canvas.Raster {
	bg := color.RGBA{0, 0, 0, 255}
	return canvas.NewRaster(func(w, h int) image.Image {
		frame := d.Image()
		ib := frame.Bounds()
		if w != ib.Dx() || h != ib.Dy() {
			rect := image.Rect(0, 0, w, h)
			cache := image.NewRGBA(rect)
			draw.Draw(cache, cache.Bounds(), &image.Uniform{bg}, image.ZP, draw.Src)
			for x := 0; x < w && x < ib.Dx(); x++ {
				for y := 0; y < h && y < ib.Dy(); y++ {
					if c := frame.At(x, y); c != nil {
						cache.Set(x, y, c)
					}
				}
			}
			frame = cache
		}
		return frame
	})
}
