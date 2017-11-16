package painter

import (
	"image/color"
	"math"
)

///////////////////////////////// Modification /////////////////////////////////

// Line has the brush move in a line.
func (b *Brush) Line(distance float64) *Brush {
	for ; distance >= 1; distance-- {
		b.modState(math.Cos(b.a)*1, math.Sin(b.a)*1, 0)
	}
	return b.modState(math.Cos(b.a)*distance, math.Sin(b.a)*distance, 0)
}

// Right turns the Brush 'angle' degrees to the right.
func (b *Brush) Right(angle float64) *Brush {
	return b.modState(0, 0, angle*toRad)
}

// Left turns the Brush 'angle' degrees to the left.
func (b *Brush) Left(angle float64) *Brush {
	return b.Right(-angle)
}

// Stop stops the Brush.
func (b *Brush) Stop() *Brush {
	b.moving = false
	return b
}

// Reset saves the state of the Brush and restores it after
// executing the given function, if the Brush is still moving
// at that point.
func (b *Brush) Reset(f func()) func() {
	return func() {
		v := *b
		f()
		if b.moving {
			*b = v
		}
	}
}

// ChangeColor changes the color of the brush
func (b *Brush) ChangeColor(color color.Color) *Brush {
	if b.moving {
		b.c = color
	}
	return b
}

// Lift lifts the Brush from the canvas
func (b *Brush) Lift() *Brush {
	return b.ChangeColor(color.Transparent)
}

// Lower lower the Brush to the canvas
func (b *Brush) Lower() *Brush {
	return b.ChangeColor(color.White)
}

// Repeat repeats function 'f' 'i' times.
func (b *Brush) Repeat(i int, f func()) *Brush {
	for ; i > 0; i-- {
		f()
	}
	return b
}
