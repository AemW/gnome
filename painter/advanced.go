package painter

import (
	"image/color/palette"
	"math"
	"math/rand"
)

var pSize = len(palette.Plan9)

// Random has the Painter draw random stuff.
func (w *Painter) Random() *Painter {
	if rand.Intn(10) > 6 {
		n := rand.Intn(5)
		switch n {
		case 0:
			w.Circle(float64(rand.Int63n(35)) + 5)
		case 1:
			w.Square(float64(rand.Int63n(35)) + 5)
		case 2:
			w.ChangeColor(palette.Plan9[rand.Intn(pSize)])
		case 3:
			w.Triangle(float64(rand.Int63n(35)) + 5)
		}
	} else {
		n := float64(rand.Intn(50))
		switch {
		case n < 10:
			w.Right(180)
			w.Line(n * 3)
		case n < 20:
			w.Line(n * 2)
		case n < 40:
			w.Right(math.Remainder(n, 10) * 10)
		case n < 50:
			w.Left(math.Remainder(n, 10) * 10)
		}
	}
	return w
}

// RandomColor changes the color of the walker's trace randomly.
func (w *Painter) RandomColor() *Painter {
	return w.ChangeColor(palette.Plan9[rand.Intn(pSize)])
}

//////////////////////////////////// --- ////////////////////////////////////

// GetHelp tells the Team to send another Painter to help with drawing
// the given function.
func (w *Painter) GetHelp(s Sketch) *Painter {
	w.phone <- func(v *Painter) {
		v.a, v.x, v.y = w.a, w.x, w.y
		v.c = w.c
		s(v)
	}
	return w
}

//////////////////////////////////// Shapes ////////////////////////////////////

const granularity = 20

// Circle has the walker move around in a circle with radius 'radius'.
func (w *Painter) Circle(radius float64) *Painter {
	stepLen := (2 * math.Pi * radius) / float64(granularity)
	return w.Polygon(granularity, stepLen)
}

// Square has the walker move around in a square with side length 'sideLen'.
func (w *Painter) Square(sideLen float64) *Painter {
	return w.Polygon(4, sideLen)
}

// Triangle has the walker move around in a triangle with side length 'sideLen'.
func (w *Painter) Triangle(sideLen float64) *Painter {
	return w.Polygon(3, sideLen)
}

// Shape traces a polygon with 'sides' sides and side length of stepLen.
// At every side the function f is called, which may be another tracing
// function. This can result in rather random shapes, since if f isn't
// created properly then the trace might not fomr an area.
func (w *Painter) Shape(sides int, stepLen float64, f func()) *Painter {
	angle := float64(360 / sides) //180 - float64(((sides-2)*180)/sides)
	turn := func() {
		w.Right(angle)
		w.Line(stepLen)
		f()
	}
	return w.Repeat(sides, turn)
}

// Polygon traces a polygon with 'sides' sides and side length of stepLen.
func (w *Painter) Polygon(sides int, stepLen float64) *Painter {
	return w.Shape(sides, stepLen, func() { /* Nothing */ })
}
