package backend

import (
	"crypto/rand"
	"image/color"
	"image/color/palette"
	"math/big"
)

// Canvas specifies the necessary implementation for a canvas
type Canvas interface {
	// Set paints a pixel with given coordinates with a given color
	Set(x, y float64, color color.Color)
	// StartEventhandler start a handler for input events
	StartEventhandler(chan bool)
	// Flush flushes (draws) the current state of the canvas
	Flush()
	// Close terminates the canvas
	Close()
	// Prepare prepares the canvas for painting and rendering
	Prepare()
}

// CanvasFactory Specifies the necessary methods for a factory for Canvas
type CanvasFactory interface {
	// Make creates a new Canvas
	Make(xSize, ySize int, title string) Canvas
	// Run kicks of any backend processes necessary for the canvas
	Run()
}

var pSize = len(palette.Plan9)

// RandomColor returns a random color.
func RandomColor() color.Color {
	n, err := randn(int64(pSize))
	if err != nil {
		return color.White
	}
	return palette.Plan9[n.Int64()]
}

func randn(tot int64) (*big.Int, error) {
	return rand.Int(rand.Reader, big.NewInt(tot))
}
