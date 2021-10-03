package handlers

import (
	"github.com/mitchellh/mapstructure"
	"github.com/streadway/amqp"
	"log"
	"proj/rabbitMq/commands"
	"reflect"
)

type Event struct {
	Command  map[string]interface{}
	Delivery amqp.Delivery
}

type IRabbitHandler interface {
	GetId() string
	GetStop() *chan struct{}
	Unsubscribe()
	Handle(event Event)
}

type RabbitHandler struct {
	Id   string
	Stop chan struct{}
}

func (r *RabbitHandler) Unsubscribe() {
	r.Stop <- struct{}{}
	log.Println("Stop func")
}

func (r *RabbitHandler) GetId() string {
	return r.Id
}

func (r *RabbitHandler) GetStop() *chan struct{} {
	return &r.Stop
}

func (*RabbitHandler) Handle(event Event) {
	log.Println("Base handler, trying execute: ", event.Command)
}

type RabbitHandlerPlus struct {
	RabbitHandler
}

func (*RabbitHandlerPlus) Handle(event Event) {
	var command commands.RabbitCommandPlus
	for key, value := range event.Command {
		log.Println("Plus handler, trying execute: ", key)
		switch key {
		case reflect.TypeOf(commands.RabbitCommandPlus{}).Name():
			err := mapstructure.Decode(value, &command)
			if err != nil {
				log.Fatal("Error")
				return
			}
			log.Println(command)

			command.Execute()
			err = event.Delivery.Ack(true)
			if err != nil {
				log.Fatal("Can not Ack")
			}
		}
	}
}

type RabbitHandlerMinus struct {
	RabbitHandler
}

func (*RabbitHandlerMinus) Handle(event Event) {
	var command commands.RabbitCommandMinus
	for key, value := range event.Command {
		log.Println("Minus handler, trying execute: ", key)
		switch key {
		case reflect.TypeOf(commands.RabbitCommandMinus{}).Name():
			err := mapstructure.Decode(value, &command)
			if err != nil {
				log.Fatal("Error")
				return
			}
			log.Println(command)

			command.Execute()

			err = event.Delivery.Ack(true)
			if err != nil {
				log.Fatal("Can not Ack")
			}
		}
	}
}
