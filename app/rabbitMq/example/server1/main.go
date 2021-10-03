package main

import (
	"proj/rabbitMq/handlers"
	"proj/rabbitMq/resolver"
	"time"
)

func main() {
	var rabbitMqConfig = resolver.RabbitMqConfig{
		Username:  "guest",
		Password:  "guest",
		Host:      "localhost",
		Port:      "5672",
		QueueName: "commands",
	}

	var r = resolver.RabbitMqResolver{
		CommandsMapChan:    make(chan handlers.Event),
		RabbitHandlersChan: make(chan handlers.IRabbitHandler),
		Subscribers:        map[string]handlers.IRabbitHandler{},
	}

	r.Start(rabbitMqConfig)
	/*var handler = handlers.RabbitHandlerMinus{
		RabbitHandler: handlers.RabbitHandler{
			Id:   "93551f4d-cb6e-45a3-8fba-2459ce6a9b70",
			Stop: make(chan struct{}),
		},
	}*/

	var handler2 = handlers.RabbitHandlerPlus{
		RabbitHandler: handlers.RabbitHandler{
			Id:   "5d173ed9-04c8-4833-a525-2888a7e362e7",
			Stop: make(chan struct{}),
		},
	}

	time.Sleep(2000)
	/*r.AddHandler(&handler)*/

	r.AddHandler(&handler2)

	forever := make(chan bool)
	<-forever
}
