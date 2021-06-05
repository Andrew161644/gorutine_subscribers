package main

import (
	guuid "github.com/google/uuid"
	"log"
	"time"
)

func main() {
	var resolver = newResolver()

	var event = eventSimple{id: guuid.New().String(), msg: "Hello"}
	var event2 = eventSimple{id: guuid.New().String(), msg: "Bye"}

	var subscriber1 = &subscriber{
		id:     guuid.New().String(),
		stop:   make(chan struct{}),
		events: make(chan IEvent),
	}

	var specialSubcriber = &specialSubscriber{subscriber{
		id:     guuid.New().String(),
		stop:   make(chan struct{}),
		events: make(chan IEvent),
	}}

	resolver.mySubcriber <- specialSubcriber
	resolver.mySubcriber <- subscriber1

	resolver.myEvent <- event
	resolver.myEvent <- event2

	specialSubcriber.HandleMessage()
	subscriber1.HandleMessage()

	time.Sleep(2000)
}

type IEvent interface {
	Id() string
	Msg() string
}

type eventSimple struct {
	id  string
	msg string
}

func (ev eventSimple) Id() string {
	return ev.id
}
func (ev eventSimple) Msg() string {
	return ev.msg
}

type resolver struct {
	myEvent     chan IEvent
	mySubcriber chan ISubscriber
}

func newResolver() *resolver {
	r := &resolver{
		myEvent:     make(chan IEvent),
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
						log.Println("Event pushed to", s)
					}
				}(id, s)
			}
		}
	}
}

type ISubscriber interface {
	GetId() string
	HandleMessage()
	GetStop() <-chan struct{}
	GetEvents() chan IEvent
}

type subscriber struct {
	id     string
	stop   <-chan struct{}
	events chan IEvent
}

type specialSubscriber struct {
	subscriber
}

func (s subscriber) GetEvents() chan IEvent {
	return s.events
}

func (s subscriber) GetStop() <-chan struct{} {
	return s.stop
}

func (s subscriber) GetId() string {
	return s.id
}

func (s subscriber) HandleMessage() {
	go func() {
		for {
			select {
			case ev := <-s.events:
				log.Println("Simple subs", ev.Msg())
			}
		}
	}()
}

/// overrides method
func (s specialSubscriber) HandleMessage() {
	go func() {
		for {
			select {
			case ev := <-s.events:
				log.Println("Special subs", ev.Msg())
			}
		}
	}()
}
