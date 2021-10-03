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
	Publish(rabbitMqConfig, commandPlus)

	var commandMinus = commands.RabbitCommandMinus{
		A: 90,
		B: 80,
	}
	Publish(rabbitMqConfig, commandMinus)

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

	err = ch.ExchangeDeclare(
		"commands", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	var commandMap = make(map[string]commands.IRabbitCommand)
	commandMap[reflect.TypeOf(command).Name()] = command

	body, _ := json.Marshal(commandMap)

	err = ch.Publish(
		"commands",
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	failOnError(err, "Failed to publish a message")
}
