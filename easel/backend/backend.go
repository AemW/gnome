package backend

import (
	"image/color"
)

type Canvas interface {
	Set(x, y float64, color color.Color)
	StartEventhandler(chan bool)
	Flush()
	Close()
}

type CanvasFactory interface {
	Make(xSize, ySize int, title string) Canvas
	Run()
}
