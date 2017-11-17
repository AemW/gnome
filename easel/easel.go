package easel

import (
	"image/color"
	"image/color/palette"
	"math/rand"
	"time"

	"github.com/AemW/gnome/easel/backend"
	glfw "github.com/AemW/gnome/easel/glfw"
	wde "github.com/AemW/gnome/easel/wde"
	"github.com/AemW/gnome/process"
)

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
	frame     Frame
}

type Frame struct {
	XSize int
	YSize int
	Title string
	Delay time.Duration
}

func NewFrame(xSize int, ySize int, title string, delay time.Duration) Frame {
	return Frame{XSize: xSize, YSize: ySize, Title: title, Delay: delay}
}

// Painter is a function which given a Canvas returns a function that when
// executed sends Pixels through to the Canvas channel.
type Painter func(Canvas) func(chan int)

func Draw(fr Frame, f func(Frame)) {
	go f(fr)
	getFactory().Run()
}

func getFactory() backend.CanvasFactory {
	i := 1
	switch i {
	case 1:
		return wde.Factory{}
	case 2:
		return glfw.Factory{}
	default:
		return wde.Factory{}
	}
}

// Make creates a new Easel of size "xSize * ySize" and 'title'
// as window title.
func Make(f Frame) Easel {
	cvs := getFactory().Make(f.XSize, f.YSize, f.Title)
	return Easel{cvs, process.Make(), make(Canvas), rand.New(rand.NewSource(time.Now().Unix())), f}
}

// PrepareEasel sets up the easel, canvas, and manager painter
func (e *Easel) PrepareEasel(manager Painter) chan bool {
	done := make(chan bool)
	e.cvs.StartEventhandler(done)
	e.prepareCanvas(e.frame.Delay)

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
		for p := range e.canvas {
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
}

var pSize = len(palette.Plan9)

// RandomColor returns a random color.
func (e *Easel) RandomColor() color.Color {
	return palette.Plan9[e.rand.Intn(pSize)]
}

/*
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
*/
