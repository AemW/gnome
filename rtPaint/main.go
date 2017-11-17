package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/AemW/gnome/easel"
	"github.com/AemW/gnome/painter"
)

func main() {

	delay := flag.Int64("d", 10, "Screen flush delay (ms)")
	size := flag.Int("s", 500, "Screen size")
	flag.Parse()

	fmt.Println("The delay is: ", *delay)
	rand.Seed(time.Now().Unix())
	f := easel.NewFrame(*size, *size, "rtPaint", time.Duration(*delay))
	easel.Draw(f, program)

}

func program(f easel.Frame) {

	// Easel instantiation
	e := easel.Make(f)

	collab, manager := painter.MakeCollaboration()

	complete := e.PrepareEasel(manager)

	sketch := func(b *painter.Brush) {
		//w.TriTriangle(40)
		b.Random()
		// nice without radian calc
		//w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
	}

	collab.Paint(float64(f.XSize/2), float64(f.YSize/2), 0, sketch)

	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 0, pg))
	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 120, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))

	<-complete
	e.Finish()
}
