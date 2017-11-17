package painter

import (
	"image/color/palette"
	"math"
	"math/rand"
)

// Random has the Painter draw random stuff.
func (b *Brush) Random() *Brush {
	if rand.Intn(10) > 6 {
		n := rand.Intn(5)
		switch n {
		case 0:
			b.Circle(float64(rand.Int63n(35)) + 5)
		case 1:
			b.Square(float64(rand.Int63n(35)) + 5)
		case 2:
			b.ChangeColor(palette.Plan9[rand.Intn(pSize)])
		case 3:
			b.Triangle(float64(rand.Int63n(35)) + 5)
		}
	} else {
		n := float64(rand.Intn(50))
		switch {
		case n < 10:
			b.Right(180)
			b.Line(n * 3)
		case n < 20:
			b.Line(n * 2)
		case n < 40:
			b.Right(math.Remainder(n, 10) * 10)
		case n < 50:
			b.Left(math.Remainder(n, 10) * 10)
		}
	}
	return b
}

var pSize = len(palette.Plan9)

// RandomColor changes the color of the walker's trace randomly.
func (b *Brush) RandomColor() *Brush {
	return b.ChangeColor(palette.Plan9[rand.Intn(pSize)])
}

//////////////////////////////////// --- ////////////////////////////////////

// GetHelp tells the Team to send another Painter to help with drawing
// the given function.
func (b *Brush) GetHelp(s Sketch) *Brush {
	b.holder <- func(v *Brush) {
		v.a, v.x, v.y = b.a, b.x, b.y
		v.c = b.c
		s(v)
	}
	return b
}

//////////////////////////////////// Shapes ////////////////////////////////////

const granularity = 30

// Circle has the walker move around in a circle with radius 'radius'.
func (b *Brush) Circle(radius float64) *Brush {
	stepLen := (2 * math.Pi * radius) / float64(granularity)
	return b.Polygon(granularity, stepLen)
}

// Square has the walker move around in a square with side length 'sideLen'.
func (b *Brush) Square(sideLen float64) *Brush {
	return b.Polygon(4, sideLen)
}

// Triangle has the walker move around in a triangle with side length 'sideLen'.
func (b *Brush) Triangle(sideLen float64) *Brush {
	return b.Polygon(3, sideLen)
}

// Shape traces a polygon with 'sides' sides and side length of stepLen.
// At every side the function f is called, which may be another tracing
// function. This can result in rather random shapes, since if f isn't
// created properly then the trace might not fomr an area.
func (b *Brush) Shape(sides int, stepLen float64, f func()) *Brush {
	angle := float64(360 / sides) //180 - float64(((sides-2)*180)/sides)
	turn := func() {
		b.Right(angle)
		b.Line(stepLen)
		f()
	}
	return b.Repeat(sides, turn)
}

// Polygon traces a polygon with 'sides' sides and side length of stepLen.
func (b *Brush) Polygon(sides int, stepLen float64) *Brush {
	return b.Shape(sides, stepLen, func() { /* Nothing */ })
}
