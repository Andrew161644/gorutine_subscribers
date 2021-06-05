package main

import (
	guuid "github.com/google/uuid"
	"log"
)

func main() {
	var resolver = newResolver()
	var id = guuid.New()
	var event = eventImpl{id: id.String(), msg: "Hello"}
	var subscriber = &subscriber{
		stop:   make(chan struct{}),
		events: make(chan IEvent),
	}
	resolver.mySubcriber <- subscriber
	resolver.myEvent <- event
	var event2 = eventImpl{id: id.String(), msg: "Bay"}
	resolver.myEvent <- event2
	subscriber.HandleMessage()
	for {

	}
}

type IEvent interface {
	Id() string
	Msg() string
}

type eventImpl struct {
	id  string
	msg string
}

func (ev eventImpl) Id() string {
	return ev.id
}
func (ev eventImpl) Msg() string {
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
			id := guuid.New()
			subscribers[id.String()] = s
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
				log.Println(ev.Msg())
			}
		}
	}()
}
