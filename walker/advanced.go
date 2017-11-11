package walker

import (
    "fmt"
    "math"
    "math/rand"
)

// Panic has the Walker perfoming a random action.
func (w *Walker) Panic() {
	n := float64(rand.Intn(33))
	switch {
	case n == 31:
		w.Circle(15)
		fmt.Println("Walker circle")
	case n == 32:
		w.Square(10)
		fmt.Println("Walker square")
	case n < 10:
		w.Walk(n)
	case n < 20:
		w.Right(math.Remainder(n, 10) * 10)
	case n < 30:
		w.Left(math.Remainder(n, 10) * 10)
	}
}

const granularity = 20

// Circle has the walker move around in a circle with radius 'radius'.
func (w *Walker) Circle(radius float64) *Walker {
	stepLen := (2 * math.Pi * radius) / float64(granularity)
	//stepLen/float64(radius)
	return w.shape(granularity, stepLen, 4*radsPerCircle)
}

// Square has the walker move around in a sqaure with side length 'sideLen'.
func (w *Walker) Square(sideLen float64) *Walker {
    return w.shape(4, sideLen, 90)
}

func (w *Walker) shape(reps int, stepLen, angle float64) *Walker {
	turn := func() {
		w.Right(angle)
		w.Walk(stepLen)
	}
	return w.Repeat(reps, turn)

}
