package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"proj/rabbitMq/commands"
	"proj/rabbitMq/resolver"
	"reflect"
)

func main() {
	var rabbitMqConfig = resolver.RabbitMqConfig{
		Username:  "guest",
		Password:  "guest",
		Host:      "localhost",
		Port:      "5672",
		QueueName: "commands",
	}

	var commandPlus = commands.RabbitCommandPlus{
		A: 50,
		B: 40,
	}

	var commandMinus = commands.RabbitCommandMinus{
		A: 90,
		B: 80,
	}
	Publish(rabbitMqConfig, commandMinus)
	Publish(rabbitMqConfig, commandPlus)

	/*time.Sleep(2000)

	var commandPlus2 = commands.RabbitCommandPlus{
		A: 60,
		B: 120,
	}

	var commandMinus2 = commands.RabbitCommandMinus{
		A: 150,
		B: 20,
	}

	Publish(rabbitMqConfig, commandPlus2)
	Publish(rabbitMqConfig, commandMinus2)*/
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func Publish(rabbitMqConfig resolver.RabbitMqConfig, command commands.IRabbitCommand) {
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

	var commandMap = make(map[string]commands.IRabbitCommand)
	commandMap[reflect.TypeOf(command).Name()] = command

	body, _ := json.Marshal(commandMap)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
	/*	log.Printf(" [x] Sent %s", body)*/
}
