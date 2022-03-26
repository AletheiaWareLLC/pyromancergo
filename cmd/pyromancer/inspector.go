package main

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/vm"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
	"strings"
	"time"
)

type PyromancerInspector struct {
	Window  fyne.Window
	Content *fyne.Container
	Tree    *widget.Tree
	Cycles  *widget.Entry
	Machine *vm.Machine
	stop    chan struct{}
	updates chan fyne.CanvasObject
}

func NewInspector(a fyne.App, m *vm.Machine) *PyromancerInspector {
	pi := &PyromancerInspector{
		Window:  a.NewWindow("Pyromancer Inspector"),
		Content: container.NewMax(),
		Tree: &widget.Tree{
			CreateNode: func(branch bool) fyne.CanvasObject {
				return widget.NewLabel("Template Object")
			},
			UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
				parts := strings.Split(uid, "/")
				node.(*widget.Label).SetText(parts[len(parts)-1])
			},
		},
		Cycles:  widget.NewEntry(),
		Machine: m,
		updates: make(chan fyne.CanvasObject, 1000),
	}

	go func() {
		for {
			select {
			case u, ok := <-pi.updates:
				if !ok {
					return
				}
				objs := make(map[fyne.CanvasObject]bool)
				objs[u] = true

				// Drain channel
				for draining := true; draining; {
					select {
					case u := <-pi.updates:
						objs[u] = true
					default:
						draining = false
					}
				}

				for o := range objs {
					o.Refresh()
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	navigation := pi.createNavigation()

	pi.Tree.ChildUIDs = func(uid string) (c []string) {
		c = navigation[uid]
		return
	}
	pi.Tree.IsBranch = func(uid string) (b bool) {
		_, b = navigation[uid]
		return
	}

	contents := pi.createContents()

	pi.Tree.OnSelected = func(uid string) {
		pi.Content.Objects = []fyne.CanvasObject{contents[uid]}
		pi.Content.Refresh()
	}
	pi.Tree.Select("Overview")

	pi.Cycles.PlaceHolder = "Tick"
	pi.Cycles.TextStyle = fyne.TextStyle{
		Monospace: true,
	}
	pi.Cycles.OnSubmitted = func(s string) {
		limit, err := strconv.Atoi(s)
		if err != nil {
			dialog.ShowError(err, pi.Window)
			return
		}
		pi.ClockWhile(func() bool {
			return pi.Machine.Tick < limit
		})
	}
	pi.Cycles.Wrapping = fyne.TextWrapOff

	pi.Window.SetContent(container.NewBorder(widget.NewToolbar(
		widget.NewToolbarSpacer(),
		NewToolbarCycleCount(pi.Cycles),
		widget.NewToolbarAction(theme.MediaPauseIcon(), pi.Stop),
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			pi.Stop()
			pi.Clock()
		}),
		widget.NewToolbarAction(theme.MediaFastForwardIcon(), pi.ClockContinuously),
		widget.NewToolbarSpacer(),
	), nil, pi.Tree, nil, pi.Content))
	pi.Window.SetOnClosed(func() {
		pi.Stop()
		close(pi.updates)
	})
	pi.Window.CenterOnScreen()
	pi.Window.Resize(fyne.NewSize(800, 600))
	pi.Window.Show()
	return pi
}

func (pi *PyromancerInspector) Clock() {
	pi.Machine.Clock()

	// Update UI
	pi.updates <- pi.Content
	pi.Cycles.Text = fmt.Sprintf("%d", pi.Machine.Tick)
	pi.updates <- pi.Cycles
}

func (pi *PyromancerInspector) ClockContinuously() {
	pi.ClockWhile(func() bool {
		return true
	})
}

func (pi *PyromancerInspector) ClockWhile(condition func() bool) {
	pi.Stop()
	pi.stop = make(chan struct{}, 1)
	go func() {
		defer func() {
			close(pi.stop)
			pi.stop = nil
		}()
		for condition() {
			select {
			case <-pi.stop:
				return
			default:
				pi.Clock()
				if pi.Machine.Processor.HasHalted() {
					pi.Stop()
				}
			}
		}
	}()
}

func (pi *PyromancerInspector) Stop() {
	if s := pi.stop; s != nil {
		select {
		case s <- struct{}{}:
		}
	}
}

func (pi *PyromancerInspector) createNavigation() map[string][]string {
	n := make(map[string][]string)

	// Top Level
	n[""] = []string{"Overview", "Processor", "Memory", "IO"}

	// Processor
	for i := 0; i < flamego.CoreCount; i++ {
		// Cores
		core := fmt.Sprintf("Processor/Core%d", i)
		n["Processor"] = append(n["Processor"], core)
		for j := 0; j < flamego.ContextCount; j++ {
			// Contexts
			context := fmt.Sprintf("Processor/Core%d/Context%d", i, j)
			n[core] = append(n[core], context)
			n[context] = []string{
				context + "/L1ICache",
				context + "/L1DCache",
			}
		}
		n[core] = append(n[core], core+"/L2Cache")
	}
	n["Processor"] = append(n["Processor"], "Processor/L3Cache")

	// IO
	n["IO"] = []string{"IO/Storage", "IO/Display"}

	return n
}

func (pi *PyromancerInspector) createContents() map[string]fyne.CanvasObject {
	contents := make(map[string]fyne.CanvasObject)

	cores := container.NewGridWithColumns(4)
	contents["Overview"] = container.NewBorder(
		widget.NewLabelWithStyle("Overview", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(
			withPaddedBorder(container.NewBorder(NewButton("Processor", func() {
				pi.Tree.OpenBranch("Processor")
				pi.Tree.Select("Processor")
			}), withPaddedBorder(NewButton("L3Cache", func() {
				pi.Tree.OpenBranch("Processor")
				pi.Tree.Select("Processor/L3Cache")
			})), nil, nil, cores)),
			withPaddedBorder(NewButton("Memory", func() {
				pi.Tree.Select("Memory")
			})),
			withPaddedBorder(container.NewBorder(NewButton("IO", func() {
				pi.Tree.Select("IO")
			}), nil, nil, nil, container.NewVBox(
				withPaddedBorder(NewButton("Storage", func() {
					pi.Tree.OpenBranch("IO")
					pi.Tree.Select("IO/Storage")
				})),
				withPaddedBorder(NewButton("Display", func() {
					pi.Tree.OpenBranch("IO")
					pi.Tree.Select("IO/Display")
				})),
			))),
		)),
	)

	contents["Processor"] = container.NewBorder(
		widget.NewLabelWithStyle("Processor", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		NewProcessorUI(pi.Machine.Processor),
	)
	contents["Processor/L3Cache"] = container.NewBorder(
		widget.NewLabelWithStyle("L3 Cache", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		NewCacheUI(pi.Machine.Processor.Cache().(*vm.Cache)),
	)

	for i := 0; i < flamego.CoreCount; i++ {
		ck := fmt.Sprintf("Processor/Core%d", i)
		c := pi.Machine.Processor.Core(i).(*vm.Core)
		contents[ck] = container.NewBorder(
			widget.NewLabelWithStyle(fmt.Sprintf("Core %d", i), fyne.TextAlignCenter, fyne.TextStyle{
				Monospace: true,
			}),
			nil, nil, nil,
			NewCoreUI(c),
		)
		contents[ck+"/L2Cache"] = container.NewBorder(
			widget.NewLabelWithStyle(fmt.Sprintf("Core %d L2 Cache", i), fyne.TextAlignCenter, fyne.TextStyle{
				Monospace: true,
			}),
			nil, nil, nil,
			NewCacheUI(c.Cache().(*vm.Cache)),
		)
		contexts := container.NewGridWithColumns(2)
		cores.Objects = append(cores.Objects, withPaddedBorder(container.NewBorder(NewButton(fmt.Sprintf("Core%d", i), func() {
			pi.Tree.OpenBranch("Processor")
			pi.Tree.OpenBranch(ck)
			pi.Tree.Select(ck)
		}), withPaddedBorder(NewButton("L2Cache", func() {
			pi.Tree.OpenBranch("Processor")
			pi.Tree.OpenBranch(ck)
			pi.Tree.Select(ck + "/L2Cache")
		})),
			nil, nil, contexts)))
		for j := 0; j < flamego.ContextCount; j++ {
			xk := fmt.Sprintf("Processor/Core%d/Context%d", i, j)
			x := c.Context(j).(*vm.Context)
			contents[xk] = container.NewBorder(
				widget.NewLabelWithStyle(fmt.Sprintf("Core %d Context %d", i, j), fyne.TextAlignCenter, fyne.TextStyle{
					Monospace: true,
				}),
				nil, nil, nil,
				NewContextUI(x),
			)
			contents[xk+"/L1ICache"] = container.NewBorder(
				widget.NewLabelWithStyle(fmt.Sprintf("Core %d Context %d L1 Instruction Cache", i, j), fyne.TextAlignCenter, fyne.TextStyle{
					Monospace: true,
				}),
				nil, nil, nil,
				NewCacheUI(x.InstructionCache().(*vm.Cache)),
			)
			contents[xk+"/L1DCache"] = container.NewBorder(
				widget.NewLabelWithStyle(fmt.Sprintf("Core %d Context %d L1 Data Cache", i, j), fyne.TextAlignCenter, fyne.TextStyle{
					Monospace: true,
				}),
				nil, nil, nil,
				NewCacheUI(x.DataCache().(*vm.Cache)),
			)
			contexts.Objects = append(contexts.Objects, withPaddedBorder(container.NewBorder(NewButton(fmt.Sprintf("Context%d", j), func() {
				pi.Tree.OpenBranch("Processor")
				pi.Tree.OpenBranch(ck)
				pi.Tree.Select(xk)
			}), nil, nil, nil, container.NewGridWithColumns(2,
				withPaddedBorder(NewButton("L1ICache", func() {
					pi.Tree.OpenBranch("Processor")
					pi.Tree.OpenBranch(ck)
					pi.Tree.OpenBranch(xk)
					pi.Tree.Select(xk + "/L1ICache")
				})),
				withPaddedBorder(NewButton("L1DCache", func() {
					pi.Tree.OpenBranch("Processor")
					pi.Tree.OpenBranch(ck)
					pi.Tree.OpenBranch(xk)
					pi.Tree.Select(xk + "/L1DCache")
				}))))))
		}
	}

	contents["Memory"] = container.NewBorder(
		widget.NewLabelWithStyle("Memory", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		NewMemoryUI(pi.Machine.Memory),
	)

	contents["IO"] = widget.NewLabelWithStyle("IO", fyne.TextAlignCenter, fyne.TextStyle{
		Monospace: true,
	})

	contents["IO/Storage"] = container.NewBorder(
		widget.NewLabelWithStyle("Storage", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		NewFileStorageUI(pi.Machine.Processor.Device(0).(*vm.FileStorage), pi.Window),
	)
	// /Users/stuartscott/Documents/Projects/Go/src/aletheiaware.com/flamego/kernel/kernel.bin

	contents["IO/Display"] = container.NewBorder(
		widget.NewLabelWithStyle("Display", fyne.TextAlignCenter, fyne.TextStyle{
			Monospace: true,
		}),
		nil, nil, nil,
		NewDisplayUI(pi.Machine.Processor.Device(1).(*vm.Display), pi.Window),
	)

	return contents
}

func withPaddedBorder(o fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(&canvas.Rectangle{
		FillColor:   color.RGBA{0, 0, 0, 32},
		StrokeColor: color.Black,
		StrokeWidth: 1,
	}, o)
}
