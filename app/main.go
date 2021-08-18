package main

import (
	guuid "github.com/google/uuid"
	events "proj/commands"
	"proj/commands/examples"
	"proj/handlers"
	"proj/handlers/examplex"
	"proj/resolvers"
	"proj/resolvers/local_resolver"
	"time"
)

func main() {
	var mResolver resolvers.IResolver = local_resolver.NewLocalResolver()
	go mResolver.GoBroadCastEvent()

	var event = examples.CalcCommandMinus{
		CommandBase: *examples.MakeEventSimple(guuid.New().String(), "Minus"),
		A:           40,
		B:           20,
	}

	var eventPlus = examples.CalcCommandPlus{
		CommandBase: *examples.MakeEventSimple(guuid.New().String(), "Plus"),
		A:           40,
		B:           20,
	}

	var subscriber1 = handlers.MakeHandler(
		guuid.New().String(),
		make(chan struct{}),
		make(chan events.ICommand))

	var subscriber2 = examplex.MakePlusSubscriber(
		guuid.New().String(),
		make(chan struct{}),
		make(chan events.ICommand))

	mResolver.AddSubscriber(subscriber1)

	mResolver.AddSubscriber(subscriber2)

	mResolver.AddCommand(event)

	mResolver.AddCommand(eventPlus)

	time.Sleep(200000)
}
