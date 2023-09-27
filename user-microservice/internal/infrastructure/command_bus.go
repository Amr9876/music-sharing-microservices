package infrastructure

import "music-sharing/user-microservice/internal/interfaces"

type CommandBus struct{}

var commandBus = &CommandBus{}

func (m *CommandBus) Send(command interfaces.CommandHandler) error {
	return command.Handle()
}

func GetCommandBus() *CommandBus {
	return commandBus
}
