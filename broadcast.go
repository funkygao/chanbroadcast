package chanbroadcast

type broadcast struct {
	c chan broadcast
	v interface{}
}

type Broadcaster struct {
	Listenc chan chan (chan broadcast)
	Sendc   chan<- interface{}
}

type Receiver struct {
	C chan broadcast
}

// create a new broadcaster object.
func NewBroadcaster() Broadcaster {
	listenc := make(chan (chan (chan broadcast)))
	sendc := make(chan interface{})
	go func() {
		currc := make(chan broadcast, 1)
		for {
			select {
			case v := <-sendc:
				if v == nil {
					currc <- broadcast{}
					return
				}
				c := make(chan broadcast, 1)
				b := broadcast{c: c, v: v}
				currc <- b
				currc = c
			case r := <-listenc:
				r <- currc
			}
		}
	}()
	return Broadcaster{
		Listenc: listenc,
		Sendc:   sendc,
	}
}

// subscribe to the broadcasts.
func (b Broadcaster) Subscribe() Receiver {
	c := make(chan chan broadcast, 0)
	b.Listenc <- c
	return Receiver{<-c}
}

// broadcast a value to all listeners.
func (b Broadcaster) Write(v interface{}) {
	b.Sendc <- v
}

// read a value that has been broadcast,
// waiting until one is available if necessary.
func (r *Receiver) Read() interface{} {
	b := <-r.C
	v := b.v
	r.C <- b
	r.C = b.c
	return v
}
