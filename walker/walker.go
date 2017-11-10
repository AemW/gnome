package walker

import (
    "math"
    "math/rand"
    "fmt"
    "github.com/AemW/screen"
    )

//import "fmt"

func (w *Walker) Panic(){
    n := float64(rand.Intn(33))
    switch {
    case n == 31:
        w.Circle(15)
        fmt.Println("Walker circle")
    case n == 32:
        w.Square(10)
        fmt.Println("Walker square")
    case n < 10:
        w.Walk(n*5)
    case n < 20:
        w.Right(math.Remainder(n,3))
    case n < 30:
        w.Left(math.Remainder(n,3))
    }
}

type Walker struct {
    x, y, a float64
    painter chan screen.Pixel
}

func AriseAt(x, y, a float64, painter screen.Painter) Walker {
    //painter <- Pixel{x, y}
    return Walker{x, y, a, painter}
}

func Arise(painter screen.Painter) Walker {
    return AriseAt(0, 0, 0, painter)
}

func WalkerProgram(x, y, a float64, f func(*Walker)) screen.Program {
    return func(painter screen.Painter) func() {
        w := AriseAt(x, y, a, painter)
        return (func() { f(&w) })
    }
}

func (w *Walker) modState(x,y,a float64) *Walker {
    w.x += x
    w.y += y
    w.a = math.Remainder(w.a + a, 360)
    w.painter <- screen.Pixel{w.x, w.y}
    return w
}

func (w *Walker) Walk(distance float64) *Walker {
    for ; distance >= 1 ; distance-- {
        w.modState(math.Cos(w.a) * 1, math.Sin(w.a) * 1, 0)
    }
    return w.modState(math.Cos(w.a) * distance, math.Sin(w.a) * distance, 0)
}


func (w *Walker) Right(angle float64) *Walker {
    return w.modState(0, 0, angle)
}

func (w *Walker) Left(angle float64) *Walker {
    return w.Right(-angle)
}

func (w *Walker) Repeat(i int, f func()) *Walker {
    for ;  i > 0; i-- { f() }
    return w
}

var granularity int = 20
func (w *Walker) Circle(radius float64) *Walker {
    step_len := (2 * math.Pi * radius)/float64(granularity)

    return w.shape(granularity, step_len, step_len/float64(radius))
}

func (w *Walker) shape(reps int, step_len, angle float64) *Walker {
    turn := func(){
        w.Right(angle)
        w.Walk(step_len)
    }
    return w.Repeat(reps, turn)

}

func (w *Walker) Square(side_len float64) *Walker{
    //return w.shape(4, side_len, side_len/2)
    return w.Repeat(4, func(){
        w.Walk(side_len)
        w.Right(90)
    })
}
