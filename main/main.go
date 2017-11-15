package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/AemW/gnome/canvas"
	"github.com/AemW/gnome/painter"
)

func main() {

	delay := flag.Int64("d", 10, "Screen flush delay (ms)")
	size := flag.Int("s", 500, "Screen size")
	flag.Parse()

	fmt.Println("The delay is: ", *delay)
	rand.Seed(time.Now().Unix())
	go program(time.Duration(*delay), *size)
	canvas.Init()

}

func program(delay time.Duration, size int) {

	// Window instantiation
	sc := canvas.Make(size, size, "Walker")

	// Start listening for input events
	done := make(chan bool)
	sc.StartEventhandler(done)

	// Painter
	sc.PrepareBrush(delay)

	h, proc := painter.Gather()
	sc.SpawnProc(proc)

	// The painter programs
	pg := func(w *painter.Painter) {
		//w.TriTriangle(40)
		w.Random()
		// nice without radian calc
		//w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
	}

	h.Paint(float64(size/2), float64(size/2), 0, pg)

	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 0, pg))
	//sc.Spawn(h.Scheme(float64(size/2), float64(size/2), 120, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))

	<-done
	sc.Stop()
	fmt.Println("canvas stopped")
}
