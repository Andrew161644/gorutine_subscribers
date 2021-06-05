package main

import (
	"strconv"
)

type IEvent interface {
	Id() string
	Msg() string
}

type eventSimple struct {
	id  string
	msg string
}

func (ev eventSimple) Id() string {
	return ev.id
}
func (ev eventSimple) Msg() string {
	return ev.msg
}

type ICalcEvent interface {
	IEvent
	Calc() int
}

type CalcEventPlus struct {
	eventSimple
	A int
	B int
}

func (ev CalcEventPlus) Msg() string {
	return strconv.Itoa(ev.Calc())
}

func (calcEvent CalcEventPlus) Calc() int {
	return calcEvent.A + calcEvent.B
}

type CalcEventMinus struct {
	eventSimple
	A int
	B int
}

func (calcEvent CalcEventMinus) Calc() int {
	return calcEvent.A - calcEvent.B
}

func (ev CalcEventMinus) Msg() string {
	return strconv.Itoa(ev.Calc())
}
