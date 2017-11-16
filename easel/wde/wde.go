package wde

import (
	"image/color"
	"math"

	wde "github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/xgb"
)

type Factory struct{}

func (f Factory) Init() {
	wde.Run()

}
func (f Factory) Make(xSize, ySize int, title string) Canvas {
	dw, err := wde.NewWindow(xSize, ySize)
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	dw.SetTitle("Title!")
	dw.SetSize(xSize, ySize)
	dw.Show()
	return Canvas{dw}
}

type Canvas struct {
	window wde.Window
}

func (c *Canvas) Set(x, y float64, color color.Color) {
	im := c.window.Screen()
	im.Set(round(x), round(y), color)
}
func (c *Canvas) Flush() {
	c.window.FlushImage()
}
func (c *Canvas) Init() {

}

func round(f float64) int {
	if f < 0 {
		return int(math.Ceil(f - 0.5))
	}
	return int(math.Floor(f + 0.5))
}
