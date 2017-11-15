package painter

import (
	"fmt"
	"image/color"
	"math"

	"github.com/AemW/gnome/canvas"
	"github.com/AemW/gnome/process"
)

// Team is a structure keeping track of all walkers.
type Team struct {
	process.Proc
	chat  Chat
	brush canvas.Brush
}

type Sketch func(*Painter)
type link chan Sketch

// Chat is a set of links
type Chat []link

// Painter represents an entiry which "walks" around the plane, sending
// it's coordinates through a channel during each step.
type Painter struct {
	x, y, a  float64
	c        color.Color
	painting bool
	brush    canvas.Brush
	phone    link
}

/////////////////////////////////// Creation ///////////////////////////////////

func makeTeam() Team {
	return Team{process.Make(), make(Chat, 0, 10), nil}
}

// AriseAt creates a new Walker at point {'x', 'y'} with an initial angle 'a'.
func (h *Team) AriseAt(x, y, a float64) Painter {
	if h.brush == nil {

	}
	//brush <- Pixel{x, y}
	l := make(link)
	h.add(l)
	return Painter{x, y, a, color.White, true, h.brush, l}
}

// Arise creates a new Walker at point {0, 0} with an initial angle 0.
func (h *Team) Arise() Painter {
	return h.AriseAt(0, 0, 0)
}

func (h *Team) Paint(x, y, a float64, f Sketch) {
	w := h.AriseAt(x, y, a)
	if h.brush == nil {
		fmt.Println("men vad faaaaan!")
	}
	h.SpawnNamed("Walker", func() {
		f(&w)
	})
}

func Gather() (*Team, canvas.Process) {
	h := makeTeam()
	return &h, func(b canvas.Brush) func(chan int) {
		h.brush = b
		return func(s chan int) { h.listen(s) }
	}
}

func (h *Team) listen(s chan int) {
	for {
		select {
		case <-s:
			h.StopPainting()
			return
		default:
			for _, c := range h.chat {
				select {
				case f := <-c:
					h.Paint(0, 0, 0, f)
				default:
					// Nothing
				}
			}
		}
	}
}

/////////////////////////////////// Commands ///////////////////////////////////

// StopPainting signals all painters to stop
func (h *Team) StopPainting() {
	h.send1(func(w *Painter) { w.Stop() })
	h.Stop()
}

func (h *Team) send(f func()) {
	h.send1(func(w *Painter) { f() })
}

func (h *Team) send1(f func(w *Painter)) {
	for _, n := range h.chat {
		n <- f
	}
}

/////////////////////////////////// Private ///////////////////////////////////

const toRad = math.Pi / 180
const radsPerCircle = math.Pi * 2

// Modify the state of the painter's brush
func (w *Painter) modState(x, y, a float64) *Painter {
	w.listen()
	if w.painting {
		w.x += x
		w.y += y
		w.a = math.Remainder(w.a+a, radsPerCircle)
		if w.brush != nil {
			w.brush <- canvas.Pixel{X: w.x, Y: w.y, Color: w.c}
		}
	}
	return w
}

func (w *Painter) listen() {
	select {
	case f := <-w.phone:
		f(w)
	default:
	}
}

func (h *Team) add(l link) {
	h.chat = append(h.chat, l)
}
