package handlers

import (
	"log"
	"proj/commands"
)

type IHandler interface {
	GetId() string
	Handle()
	GetStop() <-chan struct{}
	GetEvents() chan commands.ICommand
}

type Handler struct {
	id     string
	stop   <-chan struct{}
	Events chan commands.ICommand
}

func (s Handler) GetEvents() chan commands.ICommand {
	return s.Events
}

func (s Handler) GetStop() <-chan struct{} {
	return s.stop
}

func (s Handler) GetId() string {
	return s.id
}

func (s Handler) Handle() {
	go func() {
		for {
			select {
			case ev := <-s.Events:
				log.Println("Simple subscriber handle CommandsMapChan: ", ev.Execute())
			}
		}
	}()
}

func MakeHandler(
	id string,
	stop <-chan struct{},
	events chan commands.ICommand) *Handler {
	return &Handler{
		id:     id,
		stop:   stop,
		Events: events,
	}
}
