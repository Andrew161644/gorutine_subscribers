package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"proj/rabbitMq/handlers"
)

type RabbitMqResolver struct {
	CommandsMapChan    chan handlers.Event
	Unsubscribe        chan string
	RabbitHandlersChan chan handlers.IRabbitHandler
	Connection         *amqp.Connection
	Delivery           *<-chan amqp.Delivery
	Channel            *amqp.Channel
	Subscribers        map[string]handlers.IRabbitHandler
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (r *RabbitMqResolver) Close() {
	r.Connection.Close()
	r.Channel.Close()
}
func (r *RabbitMqResolver) Start(rabbitMqConfig RabbitMqConfig) {

	// "amqp://guest:guest@localhost:5672/"
	var connectionString = fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitMqConfig.Username,
		rabbitMqConfig.Password,
		rabbitMqConfig.Host,
		rabbitMqConfig.Port)

	var err error
	r.Connection, err = amqp.Dial(connectionString)

	failOnError(err, "Failed to connect to RabbitMQ")

	r.Channel, err = r.Connection.Channel()

	failOnError(err, "Failed to open a channel")

	err = r.Channel.ExchangeDeclare(
		"commands", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := r.Channel.QueueDeclare(
		rabbitMqConfig.QueueName, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = r.Channel.QueueBind(
		q.Name,     // queue name
		"",         // routing key
		"commands", // exchange
		false,
		nil,
	)

	var d <-chan amqp.Delivery
	d, err = r.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	r.Delivery = &d
	failOnError(err, "Failed to register a consumer")
	go r.GoBroadCastEvent()
	go r.manageCommands()
}

func (r *RabbitMqResolver) AddHandler(iSubscriber handlers.IRabbitHandler) {
	r.RabbitHandlersChan <- iSubscriber
}

func (r *RabbitMqResolver) GoBroadCastEvent() {
	for {
		select {
		case id := <-r.Unsubscribe:
			delete(r.Subscribers, id)
		case s := <-r.RabbitHandlersChan:
			r.Subscribers[s.GetId()] = s
		case commandMap := <-r.CommandsMapChan:
			for id, s := range r.Subscribers {
				go func(id string, s handlers.IRabbitHandler) {
					select {
					case <-*s.GetStop():
						r.Unsubscribe <- id
						return
					default:
						go s.Handle(commandMap)
					}
				}(id, s)
			}
		}
	}
}

func (r *RabbitMqResolver) AddCommand(event handlers.Event) {
	r.CommandsMapChan <- event
}

func (r *RabbitMqResolver) manageCommands() {
	for d := range *r.Delivery {

		var commandsMap map[string]interface{}
		json.Unmarshal(d.Body, &commandsMap)
		log.Println("Rabbit got: ", commandsMap)

		r.AddCommand(handlers.Event{
			Command:  commandsMap,
			Delivery: d,
		})
		//d.Ack(true)
	}
}
