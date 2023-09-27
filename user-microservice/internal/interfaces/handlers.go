package interfaces

type CommandHandler interface {
	Handle() error
}

type QueryHandler interface {
	Handle() (interface{}, error)
}
