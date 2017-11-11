package walker

import (
	"image/color/palette"
	"math"
	"math/rand"
)

var pSize = len(palette.Plan9)

// Panic has the Walker perfoming a random action.
func (w *Walker) Panic() {

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
			w.Walk(n * 3)
		case n < 20:
			w.Walk(n * 2)
		case n < 40:
			w.Right(math.Remainder(n, 10) * 10)
		case n < 50:
			w.Left(math.Remainder(n, 10) * 10)
		}
	}
}

const granularity = 20

//////////////////////////////////// Shapes ////////////////////////////////////

// Circle has the walker move around in a circle with radius 'radius'.
func (w *Walker) Circle(radius float64) *Walker {
	stepLen := (2 * math.Pi * radius) / float64(granularity)
	return w.shape(granularity, stepLen)
}

// Square has the walker move around in a square with side length 'sideLen'.
func (w *Walker) Square(sideLen float64) *Walker {
	return w.shape(4, sideLen)
}

// Triangle has the walker move around in a triangle with side length 'sideLen'.
func (w *Walker) Triangle(sideLen float64) *Walker {
	return w.shape(3, sideLen)
}

/*
func (w *Walker) TriTriangle(sideLen float64) *Walker {
	return w.shapeImpl(3, sideLen, func() {
		a := w.a
		w.Right(180)
		w.Triangle(sideLen / 3)
		w.a = a
	})
}
*/

func (w *Walker) shapeImpl(sides int, stepLen float64, f func()) *Walker {
	angle := float64(360 / sides) //180 - float64(((sides-2)*180)/sides)
	turn := func() {
		w.Right(angle)
		w.Walk(stepLen)
		f()
	}
	return w.Repeat(sides, turn)
}

func (w *Walker) shape(sides int, stepLen float64) *Walker {
	return w.shapeImpl(sides, stepLen, func() { /* Nothing */ })
}
