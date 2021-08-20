package commands

type ICommand interface {
	Id() string
	Execute() interface{}
}
