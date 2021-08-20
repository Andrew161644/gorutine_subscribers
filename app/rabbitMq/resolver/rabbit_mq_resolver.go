package resolver

import (
	"proj/commands"
	"proj/handlers"
	"proj/resolvers"
)

type IRabbitMqResolver interface {
	resolvers.IResolver
	Start()
	Stop()
}

type RabbitMqResolver struct {
}

func (*RabbitMqResolver) Start() {

}

func (*RabbitMqResolver) Stop() {

}

func (r *RabbitMqResolver) AddSubscriber(iSubscriber handlers.ISubscriber) {

}
func (r *RabbitMqResolver) AddCommand(command commands.ICommand) {}
func (r *RabbitMqResolver) GoBroadCastEvent()                    {}
