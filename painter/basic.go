package painter

import (
	"image/color"
	"math"
)

///////////////////////////////// Modification /////////////////////////////////

func (w *Painter) Line(distance float64) *Painter {
	for ; distance >= 1; distance-- {
		w.modState(math.Cos(w.a)*1, math.Sin(w.a)*1, 0)
	}
	return w.modState(math.Cos(w.a)*distance, math.Sin(w.a)*distance, 0)
}

// Right turns the walker 'angle' degrees to the right.
func (w *Painter) Right(angle float64) *Painter {
	return w.modState(0, 0, angle*toRad)
}

// Left turns the walker 'angle' degrees to the left.
func (w *Painter) Left(angle float64) *Painter {
	return w.Right(-angle)
}

func (w *Painter) Stop() *Painter {
	w.painting = false
	return w
}

// Reset saves the state of a painter and restores it after
// executing the given function, if the painter is still painting
// at that point.
func (w *Painter) Reset(f func()) func() {
	return func() {
		v := *w
		f()
		if w.painting {
			*w = v
		}
	}
}

// ChangeColor changes the color of the brush's trace
func (w *Painter) ChangeColor(c color.Color) *Painter {
	if w.painting {
		w.c = c
	}
	return w
}

// Lift tells the painter to lift brush from the canvas
func (w *Painter) Lift() *Painter {
	return w.ChangeColor(color.Transparent)
}

// Lower lower the painters brush to the canvas
func (w *Painter) Lower() *Painter {
	return w.ChangeColor(color.White)
}

// Repeat repeats function 'f' 'i' times.
func (w *Painter) Repeat(i int, f func()) *Painter {
	for ; i > 0; i-- {
		f()
	}
	return w
}
