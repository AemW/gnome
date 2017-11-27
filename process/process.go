package process

import (
	"fmt"
	"sync"
)

// Proc contains a mapping of all created routines.
type Proc struct {
	//stop chan bool
	wg sync.WaitGroup
	ps map[int](chan int)
	i  int
}

// Make creates a new Proc.
func Make() Proc {
	return Proc{sync.WaitGroup{}, make(map[int](chan int)), 0}
}

func (p *Proc) done() {
	fmt.Println("defered WG done")
	p.wg.Done()
}

func (p *Proc) start(stop chan int, name interface{}, f func(chan int)) {
	go func() {
		defer p.done()
		f(stop)
		fmt.Println("Stoping process: ", name)
		close(stop)
	}()
}

// SpawnListener spwans a routine with a function which listens
// to a channel.
func (p *Proc) SpawnListener(name interface{}, f func(chan int)) {
	c := make(chan int, 1)
	p.ps[p.i] = c
	p.start(c, name, f)
	fmt.Println("Starting process ", p.i, " named ", name)
	p.i++
	p.wg.Add(1)
}

// Spawn creates a new routine running function 'f'.
func (p *Proc) Spawn(f func()) {
	p.SpawnNamed(p.i, f)
}

// SpawnNamed creates a new routine named 'name' running function 'f'.
func (p *Proc) SpawnNamed(name interface{}, f func()) {
	fs := func(stop chan int) {
		for {
			select {
			default:
				f()
			case <-stop:
				return
			}
		}
	}
	p.SpawnListener(name, fs)

}

// Stop terminates all running routines.
func (p *Proc) Stop() {
	for i, c := range p.ps {
		c <- i
	}
	p.wg.Wait()
}

/*
func (p *Proc) Enumerate() {
	fmt.Println("#Processes: ", p.i)
	for i := range p.ps {
		fmt.Println("Process_id ", i)
	}
}
*/
