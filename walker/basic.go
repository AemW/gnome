package walker

import (
	"math"

	"github.com/AemW/gnome/screen"
)

// Walker represents an entiry which "walks" around the plane, sending
// it's coordinates through a channel during each step.
type Walker struct {
	x, y, a float64
	painter chan screen.Pixel
}

// AriseAt creates a new Walker at point {'x', 'y'} with an initial angle 'a' and
// a channel 'painter'.
func AriseAt(x, y, a float64, painter screen.Painter) Walker {
	//painter <- Pixel{x, y}
	return Walker{x, y, a, painter}
}

// Arise creates a new Walker at point {0, 0} with an initial angle 0 and
// a channel 'painter'.
func Arise(painter screen.Painter) Walker {
	return AriseAt(0, 0, 0, painter)
}

// Program creates a Program from a new Walker created at point {'x', 'y'}
// with an initial  angle 'a', a channel 'painter', and a function
// which modifies the walker.
func Program(x, y, a float64, f func(*Walker)) screen.Program {
	return func(painter screen.Painter) func() {
		w := AriseAt(x, y, a, painter)
		return (func() { f(&w) })
	}
}

const toRad = math.Pi / 180
const radsPerCircle = math.Pi * 2

// Modify the state of the walker
func (w *Walker) modState(x, y, a float64) *Walker {
	w.x += x
	w.y += y
	w.a = math.Remainder(w.a+a, radsPerCircle)
	w.painter <- screen.Pixel{X: w.x, Y: w.y}
	return w
}

// Walk moves the walker forward 'distance', in the direction it is facing.
func (w *Walker) Walk(distance float64) *Walker {
	for ; distance >= 1; distance-- {
		w.modState(math.Cos(w.a)*1, math.Sin(w.a)*1, 0)
	}
	return w.modState(math.Cos(w.a)*distance, math.Sin(w.a)*distance, 0)
}

// Right turns the walker 'angle' degrees to the right.
func (w *Walker) Right(angle float64) *Walker {
	return w.modState(0, 0, angle*toRad)
}

// Left turns the walker 'angle' degrees to the left.
func (w *Walker) Left(angle float64) *Walker {
	return w.Right(-angle)
}

// Repeat repeats function 'f' 'i' times.
func (w *Walker) Repeat(i int, f func()) *Walker {
	for ; i > 0; i-- {
		f()
	}
	return w
}
