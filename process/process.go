package process

import "fmt"

type Proc struct {
	stop chan bool
	ps   map[int](chan int)
	i    int
}

func Make() Proc {
	return Proc{make(chan bool), make(map[int](chan int)), 0}
}

func start(stop chan int, name interface{}, f func()) {
	go func() {
		for {
			select {
			default:
				f()
			case i := <-stop:
				close(stop)
				fmt.Println("Stoping process: ", name, " ", i)
				return
			}
		}
	}()
}

func (p *Proc) Spawn(f func()) {
	p.SpawnNamed(p.i, f)
}

func (p *Proc) SpawnNamed(name interface{}, f func()) {
	c := make(chan int)
	p.ps[p.i] = c
	start(c, name, f)
	fmt.Println("Starting process ", p.i, " named ", name)
	p.i += 1

}

func (p *Proc) Stop() {
	//for p.i >= 0 {p.ps["proc_id: " + fmt.Sprint(p.i++)] <- i }
	for i, c := range p.ps {
		//fmt.Println("Signaling process ", i)
		c <- i
	}
}

func (p *Proc) Enumerate() {
	fmt.Println("#Processes: ", p.i)
	for i := range p.ps {
		fmt.Println("Process_id ", i)
	}
}
