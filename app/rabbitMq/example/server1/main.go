package main

import (
	"bufio"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"os"
	"proj/rabbitMq/handlers"
	"proj/rabbitMq/resolver"
	"reflect"
	"strings"
	"time"
)

func main() {
	var rabbitMqConfig = resolver.RabbitMqConfig{
		Username:  "guest",
		Password:  "guest",
		Host:      "localhost",
		Port:      "5672",
		QueueName: "commands1",
	}

	var r = resolver.RabbitMqResolver{
		CommandsMapChan:    make(chan handlers.Event),
		Unsubscribe:        make(chan string),
		RabbitHandlersChan: make(chan handlers.IRabbitHandler),
		Subscribers:        map[string]handlers.IRabbitHandler{},
	}

	r.Start(rabbitMqConfig)
	time.Sleep(2000)
	reader := bufio.NewReader(os.Stdin)

	var command string
	for {
		command, _ = reader.ReadString('\n')
		splitted := strings.Split(command, " ")
		switch strings.TrimSpace(splitted[0]) {
		case "add":
			u, _ := uuid.NewV4()
			var handler2 = handlers.RabbitHandlerPlus{
				RabbitHandler: handlers.RabbitHandler{
					Id:   u.String(),
					Stop: make(chan struct{}),
				},
			}
			r.AddHandler(&handler2)
		case "minus":
			u, _ := uuid.NewV4()
			var handler2 = handlers.RabbitHandlerMinus{
				RabbitHandler: handlers.RabbitHandler{
					Id:   u.String(),
					Stop: make(chan struct{}),
				},
			}
			r.AddHandler(&handler2)
		case "ls":
			LogStatus(&r)

		case "del":
			var id = strings.TrimSpace(splitted[1])
			r.Unsubscribe <- id
		}
	}
}

func LogStatus(r *resolver.RabbitMqResolver) {
	for id, handler := range r.Subscribers {
		fmt.Printf("%s %s \n", id, reflect.TypeOf(handler))
	}
}
