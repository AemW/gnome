package wde

import (
	"image/color"
	"math"
	"runtime"

	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/xgb"

	"github.com/AemW/gnome/easel/backend"
)

type Factory struct{}

func (f Factory) Run() {
	wde.Run()

}
func (f Factory) Make(xSize, ySize int, title string) backend.Canvas {
	dw, err := wde.NewWindow(xSize, ySize)
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	dw.SetTitle("Title!")
	dw.SetSize(xSize, ySize)
	dw.Show()
	cs := Canvas{dw}
	return &cs
}

//////////////////////////// implements easel.canvas ///////////////////////////
type Canvas struct {
	window wde.Window
}

func (c *Canvas) Set(x, y float64, color color.Color) {
	im := c.window.Screen()
	im.Set(round(x), round(y), color)
}
func (c *Canvas) Flush() {
	c.window.FlushImage()
}

func (c *Canvas) Close() {
	wde.Stop()
}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

// startEventhandler starts a listener for keyboard and mouse events.
func (c *Canvas) StartEventhandler(done chan bool) {
	go func() {
		var color color.Color = color.White
		var x, y float64 = 0, 0
		dw := c.window
		events := dw.EventChan()
	loop:
		for ei := range events {
			runtime.Gosched()
			switch e := ei.(type) {
			case wde.MouseDownEvent:
				//fmt.Println("clicked", e.Where.X, e.Where.Y, e.Which)
				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						//						p = Pixel{X: float64(c.Where.X + i), Y: float64(c.Where.Y + j), Color: color}
						c.Set(float64(e.Where.X+i), float64(e.Where.Y+j), color)
						//						e.canvas <- p //Pixel{X: float64(c.Where.X + i), Y: float64(c.Where.Y + j), Color: color}
					}
				}
			case wde.MouseUpEvent:
				//color = c.RandomColor()
			case wde.MouseMovedEvent:
			case wde.MouseDraggedEvent:
			case wde.MouseEnteredEvent:
				//fmt.Println("mouse entered", e.Where.X, e.Where.Y)
			case wde.MouseExitedEvent:
				//fmt.Println("mouse exited", e.Where.X, e.Where.Y)
			case wde.MagnifyEvent:
				//fmt.Println("magnify", c.Where, c.Magnification)
			case wde.RotateEvent:
				//fmt.Println("rotate", c.Where, c.Rotation)
			case wde.ScrollEvent:
				//fmt.Println("scroll", c.Where, c.Delta)
			case wde.KeyDownEvent:
				// fmt.Println("KeyDownEvent", e.Glyph)
			case wde.KeyUpEvent:
				// fmt.Println("KeyUpEvent", e.Glyph)
			case wde.KeyTypedEvent:
				//fmt.Printf("typed key %v, glyph %v, chord %v\n", c.Key, c.Glyph, c.Chord)
				switch e.Glyph {
				case "1":
					dw.SetCursor(wde.NormalCursor)
				case "2":
					dw.SetCursor(wde.CrosshairCursor)
				case "3":
					dw.SetCursor(wde.GrabHoverCursor)
				case "4":
					//dw.Close()
					break loop
				case "c":
					//p.Color = c.RandomColor()
				}
				switch e.Key {
				case "up_arrow":
					y--
					c.Set(x, y, color)
					//e.canvas <- p
				case "down_arrow":
					y++
					c.Set(x, y, color)
					//e.canvas <- p
				case "right_arrow":
					x++
					c.Set(x, y, color)
					//e.canvas <- p
				case "left_arrow":
					x--
					c.Set(x, y, color)
					//e.canvas <- p
				}
			case wde.CloseEvent:
				//fmt.Println("close")
				//dw.Close()
				break loop
			case wde.ResizeEvent:
				//fmt.Println("resize", c.Width, c.Height)
			}
		}
		done <- true
	}()
}
