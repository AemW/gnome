package screen

import (
	"fmt"
	"image/color"
	"math"
	"runtime"
	"time"

	"github.com/AemW/gnome/process"
	"github.com/skelterjohn/go.wde"

	// Import necessary for graphical backend.
	_ "github.com/skelterjohn/go.wde/xgb"
)

// Pixel represents a pixel by its coordinates.
type Pixel struct {
	X, Y  float64
	Color color.Color
}

// Painter is a channel for Pixels.
type Painter chan Pixel

// Screen represents the current screen and the routines that modify it.
type Screen struct {
	window    wde.Window
	processes process.Proc
	painter   Painter
}

// Program is a function which given a Painter returns a function that when
// executed sends Pixels through the Painter channel.
type Program func(Painter, process.Signal) func()

// Init 'runs' the graphical backend (has to run in the main routine).
func Init() {
	wde.Run()
}

// Make creates a new Screen of size "xSize * ySize" and 'title'
// as window title.
func Make(xSize, ySize int, title string) Screen {
	// Window instantiation
	// TODO error?
	dw, err := wde.NewWindow(xSize, ySize)
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	dw.SetTitle("Title!")
	dw.SetSize(xSize, ySize)
	dw.Show()
	return Screen{dw, process.Make(), make(Painter)}
}

// SpawnPainter spawns a new goroutine which listens to the Screen's
// Painter channel and renders every received Pixel with a 'delay' ms delay.
func (s *Screen) SpawnPainter(delay time.Duration) {
	/*s.processes.SpawnNamed("Painter", func() {
		v := <-s.painter
		im := s.window.Screen()
		im.Set(round(v.X), round(v.Y), color.White)
		time.Sleep(time.Millisecond * delay)
		s.window.FlushImage()
	})*/
	go func() {
		for p := range s.painter {
			im := s.window.Screen()
			im.Set(round(p.X), round(p.Y), p.Color)
			time.Sleep(time.Millisecond * delay)
			s.window.FlushImage()
		}
	}()
}

// Spawn start a new goroutine which repeatedly executes the 'prog' Program
func (s *Screen) Spawn(prog Program) {
	s.processes.SpawnNamed("Program", func(sign process.Signal) { prog(s.painter, sign) })
}

// Stop sends a stop signal to all started routines and closes both the Painter
// channel and the graphical backend.
func (s *Screen) Stop() {
	//s.processes.Enumerate()
	s.processes.Stop()
	close(s.painter)
	s.window.Close()
	wde.Stop()
}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

// StartEventhandler starts a listener for keyboard and mouse events.
func (s *Screen) StartEventhandler(done chan bool) {
	go func(dw wde.Window) {
		events := dw.EventChan()
	loop:
		for ei := range events {
			runtime.Gosched()
			switch e := ei.(type) {
			case wde.MouseDownEvent:
				fmt.Println("clicked", e.Where.X, e.Where.Y, e.Which)
				// dw.Close()
				// break loop
			case wde.MouseUpEvent:
			case wde.MouseMovedEvent:
			case wde.MouseDraggedEvent:
			case wde.MouseEnteredEvent:
				fmt.Println("mouse entered", e.Where.X, e.Where.Y)
			case wde.MouseExitedEvent:
				fmt.Println("mouse exited", e.Where.X, e.Where.Y)
			case wde.MagnifyEvent:
				fmt.Println("magnify", e.Where, e.Magnification)
			case wde.RotateEvent:
				fmt.Println("rotate", e.Where, e.Rotation)
			case wde.ScrollEvent:
				fmt.Println("scroll", e.Where, e.Delta)
			case wde.KeyDownEvent:
				// fmt.Println("KeyDownEvent", e.Glyph)
			case wde.KeyUpEvent:
				// fmt.Println("KeyUpEvent", e.Glyph)
			case wde.KeyTypedEvent:
				fmt.Printf("typed key %v, glyph %v, chord %v\n", e.Key, e.Glyph, e.Chord)
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
				}
			case wde.CloseEvent:
				fmt.Println("close")
				//dw.Close()
				break loop
			case wde.ResizeEvent:
				fmt.Println("resize", e.Width, e.Height)
			}
		}
		done <- true
		fmt.Println("end of events")
	}(s.window)
}
