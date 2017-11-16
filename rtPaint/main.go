package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/AemW/gnome/easel"
	"github.com/AemW/gnome/painter"
)

func main() {

	delay := flag.Int64("d", 10, "Screen flush delay (ms)")
	size := flag.Int("s", 500, "Screen size")
	flag.Parse()

	fmt.Println("The delay is: ", *delay)
	go program(time.Duration(*delay), *size)
	easel.Init()

}

func program(delay time.Duration, size int) {

	// Easel instantiation
	e := easel.Make(size, size, "rtPaint")

	collab, manager := painter.MakeCollaboration()

	complete := e.PrepareEasel(delay, manager)

	sketch := func(b *painter.Brush) {
		//w.TriTriangle(40)
		b.Random()
		// nice without radian calc
		//w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
	}

	collab.Paint(float64(size/2), float64(size/2), 0, sketch)

	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 0, pg))
	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 120, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))

	<-complete
	e.Finish()
}
