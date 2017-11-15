package walker

import (
	"image/color"
	"math"

	"github.com/AemW/gnome/screen"
)

// Walker represents an entiry which "walks" around the plane, sending
// it's coordinates through a channel during each step.
type Walker struct {
	x, y, a float64
	c       color.Color
	alive   bool
	painter chan screen.Pixel
}

/////////////////////////////////// Creation ///////////////////////////////////

// AriseAt creates a new Walker at point {'x', 'y'} with an initial angle 'a' and
// a channel 'painter'.
func AriseAt(x, y, a float64, painter screen.Painter) Walker {
	//painter <- Pixel{x, y}
	return Walker{x, y, a, color.White, true, painter}
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

/////////////////////////////////// Private ///////////////////////////////////

const toRad = math.Pi / 180
const radsPerCircle = math.Pi * 2

// Modify the state of the walker
func (w *Walker) modState(x, y, a float64) *Walker {
	if w.alive {
		w.x += x
		w.y += y
		w.a = math.Remainder(w.a+a, radsPerCircle)
		w.painter <- screen.Pixel{X: w.x, Y: w.y, Color: w.c}
	}
	return w
}
