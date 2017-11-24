package painter

import "math"

func (b *Brush) SpiralOut() {
	b.Forever(func() {
		b.curve(25)
	}) /*
		b.Iter(10, func(i int) {

		})*/
}
func (b *Brush) nice() {
	b.Cont(0, func(i int) int {
		if i <= 0 {
			i = 0
		}
		b.Right(25).Line(8 + float64(i))
		return i + 1
	})
}

func (b *Brush) curve(radius float64) {
	l := (2 * math.Pi * radius) / 20
	angle := float64(360 / 20) // 180 - float64(((sides-2)*180)/sides)
	turn := func() {
		b.Line(l)
		b.Right(angle)
	}
	b.Repeat(5, turn)

}
