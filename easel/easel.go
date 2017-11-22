package easel

import (
	"image/color"
	"image/color/palette"
	"math/rand"
	"sync"
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
	Frame     Frame
}

// Frame represent a painting frame.
type Frame struct {
	XSize int
	YSize int
	Title string
	Delay time.Duration
}

// NewFrame creates a new painting frame (window) of given size and title
// and update delay.
func NewFrame(xSize int, ySize int, title string, delay time.Duration) Frame {
	return Frame{XSize: xSize, YSize: ySize, Title: title, Delay: delay}
}

// Painter is a function which given a Canvas returns a function that when
// executed sends Pixels through to the Canvas channel.
type Painter func(Canvas) func(chan int)

// Draw , given a Frame and a function, creates a Easel
// and runs the given function.
func Draw(fr Frame, f func(*Easel)) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e := makeEasel(fr)
		f(&e)
		e.finish()
	}()
	getFactory().Run()
	wg.Wait()

}

func getFactory() backend.CanvasFactory {
	i := 2
	switch i {
	case 1:
		return wde.Factory{}
	case 2:
		return glfw.Factory{}
	default:
		return wde.Factory{}
	}
}

// makeEasel creates a new Easel of size "xSize * ySize" and 'title'
// as window title.
func makeEasel(f Frame) Easel {
	cvs := getFactory().Make(f.XSize, f.YSize, f.Title)
	return Easel{cvs, process.Make(), make(Canvas), rand.New(rand.NewSource(time.Now().Unix())), f}
}

// PrepareEasel sets up the easel, canvas, and manager painter
func (e *Easel) PrepareEasel(manager Painter) chan bool {
	//done := make(chan bool)
	done := e.prepareCanvas(e.Frame.Delay)

	e.start(manager)

	return done

}

// prepareCanvas spawns a new goroutine which listens to the Easels's
// Painter channel and renders every received Pixel with a 'delay' ms delay.
func (e *Easel) prepareCanvas(delay time.Duration) chan bool {
	done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		e.cvs.Prepare()
		e.cvs.StartEventhandler(done)
		wg.Done()
		for p := range e.canvas {
			e.cvs.Set(p.X, p.Y, p.Color)
			time.Sleep(time.Millisecond * delay)
			e.cvs.Flush()
		}
	}()
	wg.Wait()
	return done
}

// start starts a new goroutine which executes the 'painter' Painter
func (e *Easel) start(painter Painter) {
	e.processes.SpawnListener("Painter", painter(e.canvas))
}

// finish sends a stop signal to all started routines and closes both the Canvas
// channel and the graphical backend.
func (e *Easel) finish() {
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
