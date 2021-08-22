package resolvers

import (
	"proj/commands"
	"proj/handlers"
)

type IHandlerContainer interface {
	AddHandler(iSubscriber handlers.IHandler)
}

type IPublisher interface {
	AddCommand(event commands.ICommand)
}
type IResolver interface {
	IHandlerContainer
	IPublisher
	GoBroadCastEvent()
}

type AbstractResolver struct {
	Commands    chan commands.ICommand
	Subscribers chan handlers.IHandler
}

func (resolver *AbstractResolver) AddHandler(iSubscriber handlers.IHandler) {
	resolver.Subscribers <- iSubscriber
}

func (resolver *AbstractResolver) AddCommand(event commands.ICommand) {
	resolver.Commands <- event
}

func NewResolver() *AbstractResolver {
	r := &AbstractResolver{
		Commands:    make(chan commands.ICommand),
		Subscribers: make(chan handlers.IHandler),
	}

	return r
}

func (resolver *AbstractResolver) GoBroadCastEvent() {
	panic("Not Implemented")
}
