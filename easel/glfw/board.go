package glfw

import "image/color"

type board [][]pixel

type pixel struct {
	x, y int
	c    color.Color
	img  uint32
}

var square = []float32{
	-0.5, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,

	-0.5, 0.5, 0,
	0.5, 0.5, 0,
	0.5, -0.5, 0,
}

func MakeBoard(xSize, ySize int) board {
	b := make([][]pixel, xSize, ySize)
	for i := 0; i < xSize; i++ {
		for j := 0; j < xSize; j++ {
			b[i] = append(b[i], makePixel(i, j, xSize, ySize))
		}
	}
	return b
}

func (b *board) paint(x, y int, c color.Color) {
	(*b)[x][y].c = c
}

func (b *board) Draw() {
	for _, ps := range *b {
		for _, p := range ps {
			p.draw()
		}
	}
}

func (p *pixel) draw() {
	//gl.BindVertexArray(c.drawable)
	//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square) / 3))
}

var sqlen = len(square)

func makePixel(x, y, height, width int) pixel {
	var points []float32
	copy(points, square)
	for i := 0; i < sqlen; i++ {
		switch i % 3 {
		case 0:
			points[i] = point(x, width, points[i])
		case 1:
			points[i] = point(y, height, points[i])
		default:
			continue
		}
	}

	return pixel{x: x, y: y, c: color.White, img: makeVao(points)}
}

func point(a, b int, p float32) float32 {
	size := 1 / float32(b)
	position := float32(a) * size
	if p < 0 {
		return (position * 2) - 1
	}
	return ((position + size) * 2) - 1

}

func makeVao(points []float32) uint32 {
	return 0
}
