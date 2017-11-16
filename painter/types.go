package painter

import (
	"image/color"
	"math"

	"github.com/AemW/gnome/easel"
	"github.com/AemW/gnome/process"
)

// Collaboration is a structure keeping track of all working painters.
type Collaboration struct {
	process.Proc
	chat   Chat
	canvas easel.Canvas
}

// Sketch is a plan for a painting waiting for a Brush
type Sketch func(*Brush)
type link chan Sketch

// Chat is a set of links
type Chat []link

// Brush represents an entiry which moves around the plane, sending
// it's coordinates through a channel during each step.
type Brush struct {
	x, y, a float64
	c       color.Color
	moving  bool
	canvas  easel.Canvas
	holder  link
}

/////////////////////////////////// Creation ///////////////////////////////////

func makeCollab() Collaboration {
	return Collaboration{process.Make(), make(Chat, 0, 10), nil}
}

// StartAt creates a new Brush at point {'x', 'y'} with an initial angle 'a'.
func (c *Collaboration) StartAt(x, y, a float64) Brush {
	l := make(link)
	c.add(l)
	return Brush{x, y, a, color.White, true, c.canvas, l}
}

// Start creates a new Brush at point {0, 0} with an initial angle 0.
func (c *Collaboration) Start() Brush {
	return c.StartAt(0, 0, 0)
}

// Paint has one painter start painting according to the given Sketch
// from the given Brush coordinates.
func (c *Collaboration) Paint(x, y, a float64, sketch Sketch) {
	b := c.StartAt(x, y, a)
	c.SpawnNamed("Brush", func() {
		sketch(&b)
	})
}

// MakeCollaboration start a new Collaboration around a manager Painter.
func MakeCollaboration() (*Collaboration, easel.Painter) {
	h := makeCollab()
	return &h, func(b easel.Canvas) func(chan int) {
		h.canvas = b
		return func(s chan int) { h.listen(s) }
	}
}

func (c *Collaboration) listen(s chan int) {
	for {
		select {
		case <-s:
			c.StopPainting()
			return
		default:
			for _, b := range c.chat {
				select {
				case f := <-b:
					c.Paint(0, 0, 0, f)
				default:
					// Nothing
				}
			}
		}
	}
}

/////////////////////////////////// Commands ///////////////////////////////////

// StopPainting signals all painters to stop painting.
func (c *Collaboration) StopPainting() {
	c.send1(func(w *Brush) { w.Stop() })
	c.Stop()
}

func (c *Collaboration) send(f func()) {
	c.send1(func(w *Brush) { f() })
}

func (c *Collaboration) send1(f func(w *Brush)) {
	for _, n := range c.chat {
		n <- f
	}
}

/////////////////////////////////// Private ///////////////////////////////////

const toRad = math.Pi / 180
const radsPerCircle = math.Pi * 2

// Modify the state of the painter's brush
func (b *Brush) modState(x, y, a float64) *Brush {
	b.listen()
	if b.moving {
		b.x += x
		b.y += y
		b.a = math.Remainder(b.a+a, radsPerCircle)
		if b.canvas != nil {
			b.canvas <- easel.Pixel{X: b.x, Y: b.y, Color: b.c}
		}
	}
	return b
}

func (b *Brush) listen() {
	select {
	case f := <-b.holder:
		f(b)
	default:
	}
}

func (c *Collaboration) add(l link) {
	c.chat = append(c.chat, l)
}
