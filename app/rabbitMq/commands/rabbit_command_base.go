package commands

import "log"

type IRabbitCommand interface {
	Execute() interface{}
}

type RabbitCommandBase struct {
	A int
	B int
}

func (r RabbitCommandBase) Execute() interface{} {
	log.Println("Exec command", r)
	return nil
}

type RabbitCommandPlus struct {
	A int
	B int
}

func (r RabbitCommandPlus) Execute() interface{} {
	var res = r.A + r.B
	log.Println("Exec command plus", res)
	return res
}

type RabbitCommandMinus struct {
	A int
	B int
}

func (r RabbitCommandMinus) Execute() interface{} {
	var res = r.A - r.B
	log.Println("Exec command minus", res)
	return res
}
