package handlers

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"proj/rabbitMq/commands"
)

type IRabbitHandler interface {
	GetId() string
	GetStop() chan struct{}
	Unsubscribe()
	Handle(command map[string]interface{})
}

type RabbitHandler struct {
	Id   string
	Stop chan struct{}
}

func (r *RabbitHandler) Unsubscribe() {
	go func() {
		r.GetStop() <- struct{}{}
		log.Println("Stop func")
	}()
}

func (r *RabbitHandler) GetId() string {
	return r.Id
}

func (r *RabbitHandler) GetStop() chan struct{} {
	return r.Stop
}

func (*RabbitHandler) Handle(command map[string]interface{}) {
	log.Println("Base handler, trying execute: ", command)
}

type RabbitHandlerPlus struct {
	RabbitHandler
}

func (*RabbitHandlerPlus) Handle(commandMap map[string]interface{}) {
	var command commands.RabbitCommandPlus
	for key, value := range commandMap {
		log.Println("Plus handler, trying execute: ", key)
		switch key {
		case "RabbitCommandPlus":
			err := mapstructure.Decode(value, &command)
			if err != nil {
				log.Fatal("Error")
				return
			}
			log.Println(command)
			command.Execute()
		}
	}
}

type RabbitHandlerMinus struct {
	RabbitHandler
}

func (*RabbitHandlerMinus) Handle(commandMap map[string]interface{}) {
	var command commands.RabbitCommandMinus
	for key, value := range commandMap {
		log.Println("Minus handler, trying execute: ", key)
		switch key {
		case "RabbitCommandMinus":
			err := mapstructure.Decode(value, &command)
			if err != nil {
				log.Fatal("Error")
				return
			}
			log.Println(command)
			command.Execute()
		}
	}
}
