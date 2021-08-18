package resolvers

import (
	"proj/commands"
	"proj/handlers"
)

type IResolver interface {
	AddSubscriber(iSubscriber handlers.ISubscriber)
	AddCommand(event commands.ICommand)
	GoBroadCastEvent()
}

type AbstractResolver struct {
	Commands    chan commands.ICommand
	Subscribers chan handlers.ISubscriber
}

func (resolver *AbstractResolver) AddSubscriber(iSubscriber handlers.ISubscriber) {
	resolver.Subscribers <- iSubscriber
}

func (resolver *AbstractResolver) AddCommand(event commands.ICommand) {
	resolver.Commands <- event
}

func NewResolver() *AbstractResolver {
	r := &AbstractResolver{
		Commands:    make(chan commands.ICommand),
		Subscribers: make(chan handlers.ISubscriber),
	}

	return r
}

func (resolver *AbstractResolver) GoBroadCastEvent() {
	panic("Not Implemented")
}
