package glfw

import (
	"fmt"
	"image/color"
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/AemW/gnome/easel/backend"
)

type Factory struct {
	done chan bool
}

func (f Factory) Run() {

}

func (f Factory) Make(xSize, ySize int, title string) backend.Canvas {
	cs := Canvas{x: xSize, y: ySize, title: title}
	return &cs
}

func (c *Canvas) Prepare() {

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(c.x, c.y, c.title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	c.context = &glfwContext{
		window:  win,
		board:   makeBoard(100, 100),
		program: shadedProgram()}

}

//////////////////////////// implements easel.canvas ///////////////////////////

type Canvas struct {
	context *glfwContext
	x, y    int
	title   string
}

func (c *Canvas) Set(x, y float64, color color.Color) {
	xf := round(x)
	yf := round(y)
	if xf < c.x && xf >= 0 && yf < c.y && yf >= 0 {
		c.context.board[xf][yf].c = color
	}
	//point := []float32{float32(x), float32(y)}
	//gl.BufferData(gl.ARRAY_BUFFER, size, data, usage)

	//im := c.window.Screen()
	//im.Set(round(x), round(y), color)
}

func (c *Canvas) Flush() {
	//gl.ClearColor(1.0, 0, 0.5, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(c.context.program)

	c.context.Draw()
	/*
		c.context.board[1][1].draw()
		c.context.board[5][5].draw()
		c.context.board[10][10].draw()
		c.context.board[7][7].draw()
		c.context.board[49][49].draw()
	*/

	glfw.PollEvents()
	c.context.window.SwapBuffers()

}

func (c *Canvas) Close() {
	glfw.Terminate()
	c.context.window.Destroy()
}

// startEventhandler starts a listener for keyboard and mouse events.
func (c *Canvas) StartEventhandler(done chan bool) {
	c.context.window.SetCharCallback(func(w *glfw.Window, char rune) {
		switch char {
		case '4':
			//w.Destroy
			done <- true
		}
	})
	/*
	       r := glfw.CharCallback(func(w *glfw.Window, c rune) {
	   		if c == rune(4) {
	   			done <- true
	   		}
	   	})
	*/

}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}
