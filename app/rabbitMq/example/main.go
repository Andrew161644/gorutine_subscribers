package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"proj/rabbitMq/commands"
	"proj/rabbitMq/handlers"
	"proj/rabbitMq/resolver"
	"reflect"
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
		CommandsMapChan:    make(chan map[string]interface{}),
		RabbitHandlersChan: make(chan handlers.IRabbitHandler),
		Subscribers:        map[string]handlers.IRabbitHandler{},
	}

	r.Start(rabbitMqConfig)

	var handler = handlers.RabbitHandlerMinus{
		RabbitHandler: handlers.RabbitHandler{
			Id:   "93551f4d-cb6e-45a3-8fba-2459ce6a9b70",
			Stop: make(chan struct{}),
		},
	}

	var handler2 = handlers.RabbitHandlerPlus{
		RabbitHandler: handlers.RabbitHandler{
			Id:   "5d173ed9-04c8-4833-a525-2888a7e362e7",
			Stop: make(chan struct{}),
		},
	}

	time.Sleep(2000)
	r.AddHandler(&handler)
	r.AddHandler(&handler2)

	Publish(rabbitMqConfig)

	log.Println("\n\n")
	handler2.Unsubscribe()

	Publish(rabbitMqConfig)

	forever := make(chan bool)
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Publish(rabbitMqConfig resolver.RabbitMqConfig) {
	var connectionString = fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitMqConfig.Username,
		rabbitMqConfig.Password,
		rabbitMqConfig.Host,
		rabbitMqConfig.Port)

	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitMqConfig.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare a queue")

	var commandPlus = commands.RabbitCommandPlus{
		A: 50,
		B: 40,
	}

	var commandMinus = commands.RabbitCommandMinus{
		A: 80,
		B: 100,
	}

	var commandMap = make(map[string]commands.IRabbitCommand)
	commandMap[reflect.TypeOf(commandPlus).Name()] = commands.RabbitCommandBase(commandPlus)
	commandMap[reflect.TypeOf(commandMinus).Name()] = commands.RabbitCommandBase(commandMinus)

	body, _ := json.Marshal(commandMap)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
