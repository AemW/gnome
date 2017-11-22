package glfw

import (
	"image/color"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const PIXEL_COLOR = "pixelColor\x00"

type glfwContext struct {
	window  *glfw.Window
	board   board
	program uint32
}

type board [][]*pixel

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
	0.5, -0.5, 0}

func makeBoard(xSize, ySize int) board {
	b := make([][]*pixel, xSize, xSize)
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

func (c *glfwContext) Draw() {

	uni := gl.GetUniformLocation(c.program, gl.Str(PIXEL_COLOR))
	for _, ps := range c.board {
		for _, p := range ps {
			r, g, b, a := p.c.RGBA()
			//fmt.Println(p.c.RGBA())
			gl.Uniform4f(uni, norm(r), norm(g), norm(b), norm(a))
			gl.BindVertexArray(p.img)
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
		}
	}
}

func norm(ui uint32) float32 {
	return float32(ui) / 0XFF
}

func makePixel(x, y, height, width int) *pixel {
	points := make([]float32, len(square), len(square))
	copy(points, square)
	for i := 0; i < len(points); i++ {
		switch i % 3 {
		case 0:
			points[i] = point(x, width, points[i])
		case 1:
			points[i] = point(y, height, points[i])
		default:
			continue
		}
	}

	return &pixel{x: x, y: y, c: color.White, img: makeVao(points)}
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
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
