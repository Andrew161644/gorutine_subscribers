package main

import "log"

type ISubscriber interface {
	GetId() string
	HandleMessage()
	GetStop() <-chan struct{}
	GetEvents() chan ICalcEvent
}

type subscriber struct {
	id     string
	stop   <-chan struct{}
	events chan ICalcEvent
}

type specialSubscriber struct {
	subscriber
}

func (s subscriber) GetEvents() chan ICalcEvent {
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
