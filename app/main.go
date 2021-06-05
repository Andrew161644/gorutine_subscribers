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
	subscriber.LogState()
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
	mySubcriber chan *subscriber
}

func newResolver() *resolver {
	r := &resolver{
		myEvent:     make(chan IEvent),
		mySubcriber: make(chan *subscriber),
	}

	go r.broadcastEvent()

	return r
}

func (r *resolver) broadcastEvent() {
	log.Println("Events listening started")
	subscribers := map[string]*subscriber{}
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
				go func(id string, s *subscriber) {
					select {
					case <-s.stop:
						unsubscribe <- id
						return
					default:
					}
					select {
					case <-s.stop:
						unsubscribe <- id
					case s.events <- e:
					}
				}(id, s)
			}
		}
	}
}

type subscriber struct {
	stop   <-chan struct{}
	events chan IEvent
}

func (s subscriber) LogState() {
	go func() {
		for {
			select {
			case ev := <-s.events:
				log.Println(ev.Msg())
			}
		}
	}()
}
