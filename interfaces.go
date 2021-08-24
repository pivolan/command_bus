package command_bus

type CommandSettings interface {
	CanSend(events ...Event) CommandSettings
	Listen(events ...Event) CommandSettings
}

type CommandBus interface {
	Command(command Command, event Event)
	SendEvent(event Event)
	RegisterCommand(commandName Command) CommandSettings
	Visualize() string
}

type Command func(bus CommandBus, event Event) error

type Event interface {
}
