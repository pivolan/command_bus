package command_bus

type CommandBusMock struct {
	Events   []Event
	Commands []Command
}

func (c *CommandBusMock) Command(command Command, event Event) {
	c.Commands = append(c.Commands, command)
}

func (c *CommandBusMock) SendEvent(event Event) {
	c.Events = append(c.Events, event)
}

func (c *CommandBusMock) RegisterCommand(commandName Command) CommandSettings {
	panic("implement me")
}

func (c *CommandBusMock) Visualize() string {
	panic("implement me")
}

func (c *CommandBusMock) Reset() {
	c.Events = []Event{}
	c.Commands = []Command{}
}
