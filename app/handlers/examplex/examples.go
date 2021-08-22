package examplex

import (
	"log"
	"proj/commands"
	"proj/commands/examples"
	"proj/handlers"
)

type specialHandler struct {
	handlers.Handler
}

type plusHandler struct {
	handlers.Handler
}

// HandleMessage / overrides method
func (s specialHandler) Handle() {
	go func() {
		for {
			select {
			case ev := <-s.Events:
				log.Println("Special subs", ev.Execute())
			}
		}
	}()
}

// HandleMessage / overrides method
func (s *plusHandler) Handle() {
	go func() {
		for {
			select {
			case ev := <-s.Events:
				eventPlus, ok := ev.(examples.CalcCommandPlus)
				if ok {
					log.Println("Plus subs", eventPlus.Execute())
				}
			}
		}
	}()
}

func MakePlusHandler(
	id string,
	stop <-chan struct{},
	events chan commands.ICommand) *plusHandler {
	return &plusHandler{*handlers.MakeHandler(id, stop, events)}
}
