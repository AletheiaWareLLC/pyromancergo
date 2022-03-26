package main

import (
	"aletheiaware.com/flamego"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func NewRegisterUI(read func(flamego.Register) uint64) fyne.CanvasObject {
	return widget.NewTable(
		func() (int, int) {
			return flamego.RegisterCount + 1 /*Header*/, 2
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
					text = "Index"
				default:
					text = "Value"
				}
			} else {
				row--
				text += "0x"
				switch col {
				case 0:
					n := nickname(row)
					if n != "" {
						n = ", " + n
					}
					// Index
					text += fmt.Sprintf("%x (r%d%s)", row, row, n)
				case 1:
					// Value
					text += fmt.Sprintf("%016x", read(flamego.Register(row)))
				}
			}
			item.(*widget.Label).SetText(text)
		},
	)
}

func nickname(i int) string {
	switch i {
	case 0:
		return "rZERO"
	case 1:
		return "rONE"
	case 2:
		return "rCID"
	case 3:
		return "rXID"
	case 4:
		return "rIVT"
	case 5:
		return "rPID"
	case 6:
		return "rPC"
	case 7:
		return "rPS"
	case 8:
		return "rPL"
	case 9:
		return "rSP"
	case 10:
		return "rSS"
	case 11:
		return "rSL"
	case 12:
		return "rDS"
	case 13:
		return "rDL"
	default:
		return ""
	}
}
