package walker

import (
	"image/color"
	"math"
)

///////////////////////////////// Modification /////////////////////////////////

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

// Die "kills" the walker, rendering it immobile.
func (w *Walker) Die() *Walker {
	w.alive = false
	return w
}

func (w *Walker) Reset(f func()) func() {
	return func() {
		v := *w
		f()
		w.a = v.a
		w.x = v.x
		w.y = v.y
		w.c = v.c
		w.alive = v.alive
	}
}

/*
//var colors = map[string]int{"white": 0, "red": 1, "blue": 2, "green": 3}

// GetColor return the color named as input, or color.White if there is
// no color named as such.
func GetColor(name string) color.Color {
	if i, ok := colors[name]; ok {
		return colorsIndex[i]
	}
	return color.White
}*/

// ChangeColor changes the color of  the walkers trace
func (w *Walker) ChangeColor(c color.Color) *Walker {
	if w.alive {
		w.c = c
	}
	return w
}

// Invisible turns the walker's trace transparent
func (w *Walker) Invisible() *Walker {
	return w.ChangeColor(color.Transparent)
}

// Visible turns the walker's trace white
func (w *Walker) Visible() *Walker {
	return w.ChangeColor(color.White)
}

// Repeat repeats function 'f' 'i' times.
func (w *Walker) Repeat(i int, f func()) *Walker {
	for ; i > 0; i-- {
		f()
	}
	return w
}
