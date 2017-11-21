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
	cs := Canvas{window: nil, x: xSize, y: ySize, title: title}
	return &cs
}

func (c *Canvas) Prepare() {

	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(c.x, c.y, c.title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	win.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	gl.ClearColor(1.0, 0.5, 0, 1.0)

	var b uint32
	gl.GenBuffers(1, &b)
	gl.BindBuffer(gl.ARRAY_BUFFER, b)

	c.window = win
}

//////////////////////////// implements easel.canvas ///////////////////////////
type Canvas struct {
	window *glfw.Window
	x, y   int
	title  string
}

func (c *Canvas) Set(x, y float64, color color.Color) {
	//point := []float32{float32(x), float32(y)}
	//gl.BufferData(gl.ARRAY_BUFFER, size, data, usage)

	//im := c.window.Screen()
	//im.Set(round(x), round(y), color)
}

func (c *Canvas) Flush() {
	//gl.ClearColor(1.0, 0, 0.5, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	c.window.SwapBuffers()
	glfw.PollEvents()

}

func (c *Canvas) Close() {
	glfw.Terminate()
}

// startEventhandler starts a listener for keyboard and mouse events.
func (c *Canvas) StartEventhandler(done chan bool) {
	c.window.SetCharCallback(func(w *glfw.Window, char rune) {
		switch char {
		case '4':
			w.Destroy()
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
