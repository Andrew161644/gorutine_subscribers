package examples

import (
	"proj/commands"
	"strconv"
)

func (ev CommandBase) Id() string {
	return ev.id
}
func (ev CommandBase) Execute() interface{} {
	return ev.msg
}

type ICalcCommand interface {
	commands.ICommand
	Calc() int
}

type CommandBase struct {
	id  string
	msg string
}

type CalcCommandPlus struct {
	CommandBase
	A int
	B int
}

func (ev CalcCommandPlus) Execute() interface{} {
	return strconv.Itoa(ev.Calc())
}

func (calcCommand CalcCommandPlus) Calc() int {
	return calcCommand.A + calcCommand.B
}

type CalcCommandMinus struct {
	CommandBase
	A int
	B int
}

func (calcCommand CalcCommandMinus) Calc() int {
	return calcCommand.A - calcCommand.B
}

func (ev CalcCommandMinus) Execute() interface{} {
	return strconv.Itoa(ev.Calc())
}

func MakeEventSimple(id string, msg string) *CommandBase {
	return &CommandBase{
		id:  id,
		msg: msg,
	}
}
