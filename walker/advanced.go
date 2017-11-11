package walker

import (
	"image/color/palette"
	"math"
	"math/rand"
)

var pSize = len(palette.Plan9)

// Panic has the Walker perfoming a random action.
func (w *Walker) Panic() {
	n := float64(rand.Intn(34))
	switch {
	case n == 31:
		w.Circle(30)
	case n == 32:
		w.Square(20)
	case n == 33:
		w.ChangeColor(palette.Plan9[rand.Intn(pSize)])
	case n < 10:
		w.Walk(n * 3)
	case n < 20:
		w.Right(math.Remainder(n, 10) * 10)
	case n < 30:
		w.Left(math.Remainder(n, 10) * 10)
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

func (w *Walker) shape(sides int, stepLen float64) *Walker {
	angle := float64(360 / sides) //180 - float64(((sides-2)*180)/sides)
	turn := func() {
		w.Right(angle)
		w.Walk(stepLen)
	}
	return w.Repeat(sides, turn)

}
