package screen

import (
	"fmt"
	"image/color"
	"math"
	"runtime"
	"time"

	"github.com/AemW/gnome/process"
	"github.com/skelterjohn/go.wde"

	// Import necessary for graphical backend
	_ "github.com/skelterjohn/go.wde/xgb"
)

// Pixel represents a pixel by its coordinates
type Pixel struct {
	X, Y float64
}

type Painter chan Pixel

type Screen struct {
	window    wde.Window
	processes process.Proc
	painter   Painter
}

type Program func(Painter) func()

func Init() {
	wde.Run()
}

func Make(x_size, y_size int, title string) Screen {
	// Window instantiation
	// TODO error?
	dw, err := wde.NewWindow(x_size, y_size)
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	dw.SetTitle("Title!")
	dw.SetSize(x_size, y_size)
	dw.Show()
	return Screen{dw, process.Make(), make(Painter)}
}

func (s *Screen) SpawnPainter(delay time.Duration) {
	s.processes.SpawnNamed("Painter", func() {
		select {
		case v := <-s.painter:
			im := s.window.Screen()
			im.Set(round(v.X), round(v.Y), color.White)
			time.Sleep(time.Millisecond * delay)
			s.window.FlushImage()
		default:
			// Nothing
		}

	})
}

func (s *Screen) Spawn(prog Program) {
	s.processes.SpawnNamed("Program", prog(s.painter))
}

func (s *Screen) Stop() {
	//s.processes.Enumerate()
	s.processes.Stop()
	close(s.painter)
	wde.Stop()
}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}

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
					dw.Close()
					break loop
				}
			case wde.CloseEvent:
				fmt.Println("close")
				dw.Close()
				break loop
			case wde.ResizeEvent:
				fmt.Println("resize", e.Width, e.Height)
			}
		}
		done <- true
		fmt.Println("end of events")
	}(s.window)
}
