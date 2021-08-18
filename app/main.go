package main

import (
	guuid "github.com/google/uuid"
	"log"
	"time"
)

func main() {
	var resolver = newResolver()

	var event = CalcEventMinus{
		eventSimple: eventSimple{id: guuid.New().String(), msg: "Minus"},
		A:           40,
		B:           20,
	}

	var eventPlus = CalcEventPlus{
		eventSimple: eventSimple{id: guuid.New().String(), msg: "Minus"},
		A:           40,
		B:           20,
	}

	var subscriber1 = &subscriber{
		id:     guuid.New().String(),
		stop:   make(chan struct{}),
		events: make(chan ICalcEvent),
	}

	var subscriber2 = &plusSubscriber{
		subscriber{
			id:     guuid.New().String(),
			stop:   make(chan struct{}),
			events: make(chan ICalcEvent),
		},
	}

	resolver.mySubcriber <- subscriber1

	resolver.mySubcriber <- subscriber2

	resolver.myEvent <- event

	resolver.myEvent <- eventPlus

	subscriber1.HandleMessage()

	subscriber2.HandleMessage()

	time.Sleep(2000)
}

type resolver struct {
	myEvent     chan ICalcEvent
	mySubcriber chan ISubscriber
}

func newResolver() *resolver {
	r := &resolver{
		myEvent:     make(chan ICalcEvent),
		mySubcriber: make(chan ISubscriber),
	}

	go r.broadcastEvent()

	return r
}

func (r *resolver) broadcastEvent() {
	log.Println("Events listening started")
	subscribers := map[string]ISubscriber{}
	unsubscribe := make(chan string)
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-r.mySubcriber:
			subscribers[s.GetId()] = s
		case e := <-r.myEvent:
			for id, s := range subscribers {
				go func(id string, s ISubscriber) {
					select {
					case <-s.GetStop():
						unsubscribe <- id
						return
					default:
					}
					select {
					case <-s.GetStop():
						unsubscribe <- id
					case s.GetEvents() <- e:
						log.Println("Event pushed to", s.GetId())
					}
				}(id, s)
			}
		}
	}
}
