package process

import "fmt"

// Proc contains a mapping of all created routines.
type Proc struct {
	//stop chan bool
	ps map[int](chan int)
	i  int
}

type Signal chan int

// Make creates a new Proc.
func Make() Proc {
	return Proc{make(map[int](chan int)), 0}
}

// TODO This type of setup disallows looping 'f'
// TODO change Program type to give more controll to caller?
func start(s Signal, name interface{}, f func(Signal)) {
	go func() {
		f(s)
		fmt.Println("Stoping process: ", name)
		close(s)
		/*
			for {
				select {
				default:
					f()
				case i := <-stop:
					fmt.Println("Stoping process: ", name, " ", i)
					close(stop)
					return
				}
			}*/
	}()
}

// Spawn creates a new routine running function 'f'.
func (p *Proc) Spawn(f func(Signal)) {
	p.SpawnNamed(p.i, f)
}

// SpawnNamed creates a new routine named 'name' running function 'f'.
func (p *Proc) SpawnNamed(name interface{}, f func(Signal)) {
	c := make(chan int)
	p.ps[p.i] = c
	start(c, name, f)
	fmt.Println("Starting process ", p.i, " named ", name)
	p.i++

}

// Stop terminates all running routines.
func (p *Proc) Stop() {
	for i, c := range p.ps {
		c <- i
	}
}

/*
func (p *Proc) Enumerate() {
	fmt.Println("#Processes: ", p.i)
	for i := range p.ps {
		fmt.Println("Process_id ", i)
	}
}
*/
