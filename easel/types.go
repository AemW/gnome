package easel

import (
	"image/color"
	"time"

	"github.com/AemW/gnome/easel/backend"
	"github.com/AemW/gnome/easel/glfw"
	"github.com/AemW/gnome/easel/wde"
	"github.com/AemW/gnome/process"
)

// Pixel represents a pixel by its coordinates and color.
type Pixel struct {
	X, Y  float64
	Color color.Color
}

// Canvas is a channel for Pixels.
type Canvas chan Pixel

// Easel represents the current frame, canvas and the routines that modify it.
type Easel struct {
	cvs       backend.Canvas
	processes process.Proc
	canvas    Canvas
	Frame     Frame
}

// Frame represent a painting frame.
type Frame struct {
	XSize int
	YSize int
	Title string
	Delay time.Duration
}

// Painter is a function which given a Canvas returns a function that when
// executed sends Pixels through to the Canvas channel.
type Painter func(Canvas) func(chan int)

//////////////////////////////// Backend engine ////////////////////////////////
// Engine is an enum interface
type Engine interface {
	Engine() engine
}

type engine int

var _engine = WDE

// Enum constants for the backend engine
const (
	// WDE is the engine built based on github.com/skelterjohn/go.wde
	WDE engine = iota
	// GLFW is the engine built based on github.com/go-gl/glfw/v3.2/glfw"
	GLFW
)

// Engine returns the enum implementation
func (e engine) Engine() engine {
	return e
}

// SetEngine changes the current backend engine
func SetEngine(b Engine) {
	_engine = b.Engine()
}

func getFactory() backend.CanvasFactory {
	switch _engine {
	case WDE:
		return wde.Factory{}
	case GLFW:
		return glfw.Factory{}
	default:
		return wde.Factory{}
	}
}

////////////////////////////////////////////////////////////////////////////////
