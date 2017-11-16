package easel

import (
	"image/color"
	"image/color/palette"
	"math/rand"
	"runtime"
	"time"

	backend "github.com/AemW/gnome/easel/wde"
	"github.com/AemW/gnome/process"
	"github.com/skelterjohn/go.wde"
	// Import necessary for graphical backend.
	_ "github.com/skelterjohn/go.wde/xgb"
)

type canvas interface {
	Set(p *Pixel)
	Flush()
	Init()
	Channel() Canvas
	Close()
}

type canvasFactory interface {
	Make(xSize, ySize int, title string) canvas
	Init()
}

// Pixel represents a pixel by its coordinates and color.
type Pixel struct {
	X, Y  float64
	Color color.Color
}

// Canvas is a channel for Pixels.
type Canvas chan Pixel

// Easel represents the current canvas and the routines that modify it.
type Easel struct {
	cvs       backend.Canvas
	processes process.Proc
	canvas    Canvas
	rand      *rand.Rand
}

// Painter is a function which given a Canvas returns a function that when
// executed sends Pixels through to the Canvas channel.
type Painter func(Canvas) func(chan int)

// Init 'runs' the graphical backend (has to run in the main routine).
func Init() {
	backend.Factory{}.Init()
}

// Make creates a new Easel of size "xSize * ySize" and 'title'
// as window title.
func Make(xSize, ySize int, title string) Easel {
	fac := backend.Factory{}
	return Easel{fac.Make(xSize, ySize, title), process.Make(), make(Canvas), rand.New(rand.NewSource(time.Now().Unix()))}
}

// PrepareEasel sets up the easel, canvas, and manager painter
func (e *Easel) PrepareEasel(delay time.Duration, manager Painter) chan bool {
	done := make(chan bool)
	e.startEventhandler(done)
	e.prepareCanvas(delay)

	e.start(manager)

	return done

}

// prepareCanvas spawns a new goroutine which listens to the Easels's
// Painter channel and renders every received Pixel with a 'delay' ms delay.
func (e *Easel) prepareCanvas(delay time.Duration) {
	/*s.processes.SpawnNamed("Painter", func() {
		v := <-s.brush
		im := s.window.Screen()
		im.Set(round(v.X), round(v.Y), color.White)
		time.Sleep(time.Millisecond * delay)
		s.window.FlushImage()
	})*/
	go func() {
		for p := range e.cvs.Channel() {
			e.cvs.Set(p.X, p.Y, p.Color)
			time.Sleep(time.Millisecond * delay)
			e.cvs.Flush()
		}
	}()
}

// start starts a new goroutine which executes the 'painter' Painter
func (e *Easel) start(painter Painter) {
	e.processes.SpawnListener("Painter", painter(e.canvas))
}

// Finish sends a stop signal to all started routines and closes both the Canvas
// channel and the graphical backend.
func (e *Easel) Finish() {
	//s.processes.Enumerate()
	e.processes.Stop()
	close(e.canvas)
	e.cvs.Close()
	wde.Stop()
}

var pSize = len(palette.Plan9)

// RandomColor returns a random color.
func (e *Easel) RandomColor() color.Color {
	return palette.Plan9[e.rand.Intn(pSize)]
}

// startEventhandler starts a listener for keyboard and mouse events.
func (e *Easel) startEventhandler(done chan bool) {
	go func(dw wde.Window) {
		var color color.Color = color.White
		p := Pixel{X: 0, Y: 0, Color: color}
		events := dw.EventChan()
	loop:
		for ei := range events {
			runtime.Gosched()
			switch c := ei.(type) {
			case wde.MouseDownEvent:
				//fmt.Println("clicked", e.Where.X, e.Where.Y, e.Which)
				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						p = Pixel{X: float64(c.Where.X + i), Y: float64(c.Where.Y + j), Color: color}
						e.canvas <- p //Pixel{X: float64(c.Where.X + i), Y: float64(c.Where.Y + j), Color: color}
					}
				}
			case wde.MouseUpEvent:
				color = e.RandomColor()
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
				switch c.Glyph {
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
					p.Color = e.RandomColor()
				}
				switch c.Key {
				case "up_arrow":
					p.Y--
					//fmt.Println(p)
					e.canvas <- p
				case "down_arrow":
					p.Y++
					//fmt.Println(p)
					e.canvas <- p
				case "right_arrow":
					p.X++
					//fmt.Println(p)
					e.canvas <- p
				case "left_arrow":
					p.X--
					//fmt.Println(p)
					e.canvas <- p
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
	}(e.window)
}
