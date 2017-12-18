package painter

func (b *Brush) SpiralOut(init, inc int, len, rad float64) {
	b.Cont(init, func(i int) int {
		b.Curve(len+float64(i), rad+float64(i))
		return i + inc
	})
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

func (b *Brush) Weird() {
	b.Shape(15, 25, func() {
		sketch := func(b *Brush) {
			b.RandomColor().SpiralOut(15, 2, 2, 2)
		}
		b.GetHelp(sketch)
	})
	b.Stop()

}
