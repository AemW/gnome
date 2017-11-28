package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/AemW/gnome/easel"
	"github.com/AemW/gnome/painter"
)

func main() {

	delay := flag.Int64("d", 1, "Screen flush delay (ms)")
	size := flag.Int("s", 500, "Screen size")
	flag.Parse()

	log.Println("The delay is: ", *delay)
	rand.Seed(time.Now().Unix())

	easel.SetEngine(easel.GLFW)
	f := easel.NewFrame(*size, *size, "rtPaint", time.Duration(*delay))
	easel.Draw(f, program)
	log.Println("Exiting")

}

func program(e *easel.Easel) {

	collab, manager := painter.MakeCollaboration()
	//_, manager := painter.MakeCollaboration()

	complete := e.PrepareEasel(manager)

	sketch := func(b *painter.Brush) {
		//w.TriTriangle(40)
		//b.Random()
		b.Weird()
		//b.SpiralOut(10, 0, 0)
		// nice without radian calc
		//w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
	}

	collab.Paint(float64(e.Frame.XSize/2), float64(e.Frame.YSize/2), 0, sketch)

	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 0, pg))
	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 120, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))

	log.Println("Waiting for completion signal")
	<-complete
	log.Println("Recieved completion signal")

}
