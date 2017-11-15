package walker

import (
	"fmt"
	"image/color"
	"math"

	"github.com/AemW/gnome/screen"
)

type Hive struct {
	net HiveNet
}

type Command func(W *Walker)
type Link chan Command
type HiveNet []Link

// Walker represents an entiry which "walks" around the plane, sending
// it's coordinates through a channel during each step.
type Walker struct {
	x, y, a float64
	c       color.Color
	alive   bool
	painter screen.Painter
	hivenet Link
}

/////////////////////////////////// Creation ///////////////////////////////////

func MakeHive() Hive {
	return Hive{make(HiveNet, 0, 10)}
}

// AriseAt creates a new Walker at point {'x', 'y'} with an initial angle 'a' and
// a channel 'painter'.
func (h *Hive) AriseAt(x, y, a float64, painter screen.Painter) Walker {
	//painter <- Pixel{x, y}
	link := make(Link)
	h.add(link)
	return Walker{x, y, a, color.White, true, painter, link}
}

// Arise creates a new Walker at point {0, 0} with an initial angle 0 and
// a channel 'painter'.
func (h *Hive) Arise(painter screen.Painter) Walker {
	return h.AriseAt(0, 0, 0, painter)
}

// Scheme creates a Scheme from a new Walker created at point {'x', 'y'}
// with an initial  angle 'a', a channel 'painter', and a function
// which modifies the walker.
func (h *Hive) Scheme(x, y, a float64, f func(*Walker)) screen.Program {
	return func(painter screen.Painter) func() {
		w := h.AriseAt(x, y, a, painter)
		return (func() { f(&w) })
	}
}

/////////////////////////////////// Commands ///////////////////////////////////

func (h *Hive) Stop() {
	h.send1(func(w *Walker) { w.Die() })
}

func (h *Hive) SayHello() {
	h.send(func() { fmt.Println("Hellloooozzz") })
}

func (h *Hive) send(f func()) {
	h.send1(func(w *Walker) { f() })
}
func (h *Hive) send1(f func(w *Walker)) {
	for _, n := range h.net {
		n <- f
	}
}

/////////////////////////////////// Private ///////////////////////////////////

const toRad = math.Pi / 180
const radsPerCircle = math.Pi * 2

// Modify the state of the walker
func (w *Walker) modState(x, y, a float64) *Walker {
	w.listen()
	if w.alive {
		w.x += x
		w.y += y
		w.a = math.Remainder(w.a+a, radsPerCircle)
		w.painter <- screen.Pixel{X: w.x, Y: w.y, Color: w.c}
	}
	return w
}

func (w *Walker) listen() {
	select {
	case f := <-w.hivenet:
		f(w)
	default:
	}
}

func (h *Hive) add(l Link) {
	h.net = append(h.net, l)
}
