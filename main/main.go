package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/AemW/gnome/screen"
	"github.com/AemW/gnome/walker"
)

func main() {

	delay := flag.Int64("d", 10, "Screen flush delay (ms)")
	size := flag.Int("s", 500, "Screen size")
	flag.Parse()

	fmt.Println("The delay is: ", *delay)
	rand.Seed(time.Now().Unix())
	go program(time.Duration(*delay), *size)
	screen.Init()

}

func program(delay time.Duration, size int) {

	// Window instantiation
	sc := screen.Make(size, size, "Walker")

	// Start listening for input events
	done := make(chan bool)
	sc.StartEventhandler(done)

	// Painter
	sc.SpawnPainter(delay)

	// The walker programs
	pg := func(w *walker.Walker) {
		//w.TriTriangle(40)
		w.Panic()
		// nice without radian calc
		//w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
	}

	sc.Spawn(walker.Program(float64(size/2), float64(size/2), 0, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 120, pg))
	//sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))

	<-done
	sc.Stop()
	fmt.Println("screen stopped")
}
