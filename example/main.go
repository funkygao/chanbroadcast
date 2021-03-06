package main

import (
	"fmt"
	broadcast "github.com/funkygao/chanbroadcast"
	"time"
)

var b = broadcast.NewBroadcaster()

func listen(r broadcast.Receiver) {
	for v := r.Read(); v != nil; v = r.Read() {
		go listen(r)
		fmt.Println(v)
	}
}

func main() {
	r := b.Subscribe()
	go listen(r)

	// produce data
	for i := 0; i < 10; i++ {
		b.Write(i)
	}

	// EOF signal
	b.Write(nil)

	time.Sleep(3 * 1e9)
}
