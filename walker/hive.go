package walker

import (
	"github.com/AemW/gnome/process"
	"github.com/AemW/gnome/screen"
)

type Hive struct {
	ws      map[int]Walker
	painter screen.Painter
	fanIn   map[int]screen.Painter
	index   int
}

func HiveMind() Hive {
	//TODO nil
	return Hive{make(map[int]Walker), nil, make(map[int]screen.Painter), 0}
}

func (h *Hive) Program(x, y, a float64, f func(*Walker)) screen.Program {
	c := make(screen.Painter)
	h.fanIn[h.index] = c
	h.index++
	return func(painter screen.Painter, s process.Signal) func() {
		h.painter = painter
		w := AriseAt(x, y, a, c)
		h.ws[h.index] = w
		h.index++
		select {
		case v := <-s:

		default:
			h.Next()
		}
		return (func() {

			f(&w)
		})
	}
}

func (h *Hive) Next() {
	is := make([]int, h.index)
	for i, j := 0, 0; i < h.index; i++ {
		if v, ok := <-h.fanIn[i]; ok {
			h.painter <- v
		} else {
			is[j] = i
			j++
		}
	}
	for i := range is {
		delete(h.fanIn, i)
		h.index--
	}
}
