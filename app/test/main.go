package main

import (
	"fmt"
	"time"
)

type It interface {
	GetA() *chan struct{}
}
type N struct {
	A chan struct{}
}

func (n *N) GetA() *chan struct{} {
	return &n.A
}

func main() {
	var ch = make(chan struct{})
	var n It = &N{
		A: ch,
	}

	go func(ch *chan struct{}) {
		for {
			select {
			case <-*ch:
				fmt.Println("H")
			default:
				fmt.Println("N")
			}
		}
	}(&ch)
	time.Sleep(20000)
	*n.GetA() <- struct{}{}
}
