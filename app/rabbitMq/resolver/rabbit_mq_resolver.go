package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"proj/rabbitMq/handlers"
	"proj/resolvers"
)

type IRabbitMqResolver interface {
	resolvers.IHandlerContainer
	Start(RabbitMqConfig)
	Close()
}

type RabbitMqResolver struct {
	CommandsMapChan    chan map[string]interface{}
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

	q, err := r.Channel.QueueDeclare(
		rabbitMqConfig.QueueName, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	failOnError(err, "Failed to declare a queue")

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
	log.Println("Add handler: ", iSubscriber.GetId())
	r.RabbitHandlersChan <- iSubscriber
}

func (r *RabbitMqResolver) GoBroadCastEvent() {
	log.Println("examples listening started")
	unsubscribe := make(chan string)
	for {
		select {
		case id := <-unsubscribe:
			delete(r.Subscribers, id)
		case s := <-r.RabbitHandlersChan:
			r.Subscribers[s.GetId()] = s
		case commandMap := <-r.CommandsMapChan:
			for id, s := range r.Subscribers {
				go func(id string, s handlers.IRabbitHandler) {
					select {
					case <-s.GetStop():
						unsubscribe <- id
						return
					default:
						go s.Handle(commandMap)
					}
				}(id, s)
			}
		}
	}
}

func (r *RabbitMqResolver) AddCommand(event map[string]interface{}) {
	r.CommandsMapChan <- event
}

func (r *RabbitMqResolver) manageCommands() {
	for d := range *r.Delivery {

		var commandsMap map[string]interface{}
		err := json.Unmarshal(d.Body, &commandsMap)
		log.Println("Rabbit got: ", commandsMap)

		r.AddCommand(commandsMap)
		err = d.Ack(true)
		if err != nil {
			log.Fatal("Can not Ack")
		}
	}
}
