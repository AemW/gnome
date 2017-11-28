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

var (
	BUFFER_SIZE int = 4
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

	//glfw.SwapInterval(10)

	c.context = &glfwContext{
		window:  win,
		board:   makeBoard(c.x, c.y),
		program: shadedProgram()}

}

//////////////////////////// implements easel.canvas ///////////////////////////

type Canvas struct {
	context *glfwContext
	x, y    int
	title   string
}

// Set sets the color of a pixel.
func (c *Canvas) Set(x, y float64, cl color.Color) {
	xf := round(x)
	yf := round(y)
	if xf < c.x && xf >= 0 && yf < c.y && yf >= 0 {
		c.context.clrd = append(c.context.clrd, c.context.board[xf][yf])
		c.context.board[xf][yf].c = cl
		/*
			p := c.context.board[xf][yf]
			if p.c == color.Black {
				c.context.clrd = append(c.context.clrd, p)
			}
			p.c = cl
		*/
	}
}

func (c *Canvas) Flush() {
	// Not necessary to clear the screen since we want the previous
	// traces to remain
	//	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	if len(c.context.clrd) >= BUFFER_SIZE {
		gl.UseProgram(c.context.program)

		c.context.Draw()

		glfw.PollEvents()
		c.context.window.SwapBuffers()
	}

}

func (c *Canvas) Close() {
	c.context.window.Destroy()
	glfw.Terminate()
}

// StartEventhandler starts a listener for keyboard and mouse events.
func (c *Canvas) StartEventhandler(done chan bool) {
	c.context.window.SetCharCallback(func(w *glfw.Window, char rune) {
		switch char {
		case '4':
			done <- true
		case '+':
			BUFFER_SIZE++
		case '-':
			if BUFFER_SIZE > 1 {
				BUFFER_SIZE--
			}
		}
	})
}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}
