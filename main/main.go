package main

import ("fmt"
        "math/rand"
        "time"
        "flag"
        "github.com/skelterjohn/go.wde"
        "github.com/AemW/walker"
        "github.com/AemW/screen"
    )



func main() {

    delay := flag.Int64("d", 10, "Screen flush delay (ms)")
    size := flag.Int("s", 500, "Screen size")
    flag.Parse()

    fmt.Println("The delay is: ", *delay)
    rand.Seed(time.Now().Unix())
    go program(time.Duration(*delay), *size)
    wde.Run()


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
        w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
    }

    sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 0, pg))
    //sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 120, pg))
    //sc.Spawn(walker.WalkerProgram(float64(size/2), float64(size/2), 240, pg))
    /*
    p.Spawn(func(){
        //w.Panic()
        w.Walk(10).Right(90).Walk(10).Right(45).Walk(10)
    })

    p.Spawn(func(){
        //w.Panic()
        w.Right(120).Walk(10).Right(90).Walk(10).Right(45).Walk(10)
    })

    p.Spawn(func(){
        //w.Panic()
        w.Right(240).Walk(10).Right(90).Walk(10).Right(45).Walk(10)
    })
    /*
    go func () {
        for {
            select {
            default:
                w.Panic()
            case <-stop:
                return
            }
        }
    }()
    */
    <- done
    sc.Stop()
    fmt.Println("screen stopped")
    wde.Stop()
}