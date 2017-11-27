package easel

import (
	"log"
	"sync"
	"time"

	"github.com/AemW/gnome/process"
)

// NewFrame creates a new painting frame (window) of given size and title
// and update delay.
func NewFrame(xSize int, ySize int, title string, delay time.Duration) Frame {
	return Frame{XSize: xSize, YSize: ySize, Title: title, Delay: delay}
}

// Draw , given a Frame and a function, creates a Easel
// and runs the given function.
func Draw(fr Frame, f func(*Easel)) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	e := makeEasel(fr)

	// Spawn a routine that instantiates the main components and kicks-off
	// the rendering somehow specefied by function f
	go func() {
		defer wg.Done()
		f(&e)
		log.Println("Finishing main loop")

		log.Println("Cleaning easel")
		e.finish()

	}()

	// Backend specific stuff, which needs to be done in the background
	getFactory().Run()
	wg.Wait()

}

// makeEasel creates a new Easel with dimensions as specified by the frame.
func makeEasel(f Frame) Easel {
	cvs := getFactory().Make(f.XSize, f.YSize, f.Title)
	return Easel{cvs, process.Make(), make(Canvas), f}
}

// PrepareEasel sets up the easel, canvas, and manager painter
func (e *Easel) PrepareEasel(manager Painter) chan bool {
	done := e.prepareCanvas()

	e.start(manager)

	return done
}

// prepareCanvas spawns a new goroutine which listens to the Easels's
// Painter channel and renders every received Pixel with a delay.
func (e *Easel) prepareCanvas() chan bool {
	done := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Spawn the main graphical thread which is responsible for
	// rendering
	go func() {

		// Prepare the rendering backend
		e.cvs.Prepare()
		e.cvs.StartEventhandler(done)
		wg.Done()

		// Paint a pixel with a given color then draw it after a short delay
		for p := range e.canvas {
			e.cvs.Set(p.X, p.Y, p.Color)
			time.Sleep(time.Millisecond * e.Frame.Delay)
			e.cvs.Flush()
		}
	}()

	// Wait for the backend to finish preparing
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
	e.processes.Stop()
	close(e.canvas)
	e.cvs.Close()
}
