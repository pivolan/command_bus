package command_bus

import (
	"fmt"
	"github.com/bdlm/log"
	"reflect"
	"strings"
)

type commandName string
type eventName string
type commandBus struct {
	listeners          map[commandName][]eventName
	listenersSendEvent map[eventName][]Command
	canSend            map[commandName]map[eventName]bool
	logger             *log.Logger
}

func NewCommandBus(logger *log.Logger) CommandBus {
	return &commandBus{
		listeners:          map[commandName][]eventName{},
		listenersSendEvent: map[eventName][]Command{},
		canSend:            map[commandName]map[eventName]bool{},
		logger:             logger,
	}
}

type commandSettings struct {
	CommandSettings
	bus         *commandBus
	commandName commandName
	command     Command
}

func (c *commandSettings) CanSend(events ...Event) CommandSettings {
	for _, event := range events {
		c.bus.canSend[c.commandName][getEventName(event)] = true
	}
	return c
}

func (c *commandSettings) Listen(events ...Event) CommandSettings {
	for _, event := range events {
		c.bus.listeners[c.commandName] = append(c.bus.listeners[c.commandName], getEventName(event))
		c.bus.listenersSendEvent[getEventName(event)] = append(c.bus.listenersSendEvent[getEventName(event)], c.command)
	}
	return c
}

func (c *commandBus) Command(command Command, event Event) {
	c.logger.Debugf("%s Command, event %s, %v\n", getCommandName(command), getEventName(event), getFieldsList(event))
	go func() {
		err := command(c, event)
		if err != nil {
			c.logger.Errorf("Errors on: Command %s, event %s, err: %s", getCommandName(command), getEventName(event), err)
		}
	}()
}

func (c *commandBus) SendEvent(event Event) {
	caller := myCaller()
	c.logger.Debugf(fmt.Sprintf("%s(%v) Event sent by command %s\n", getEventName(event), getFieldsList(event), caller))
	if _, ok := c.canSend[caller]; !ok {
		c.logger.Fatalf("Command not registered in commandBus, only registered commands can send events, caller: %s, event: %s\n", caller, getEventName(event))
	}
	if _, ok := c.canSend[caller][getEventName(event)]; !ok {
		c.logger.Fatalf("This event not allowed to call by command, caller: %s, event: %s\n", caller, getEventName(event))
	}
	eName := eventName(reflect.TypeOf(event).String())
	for _, command := range c.listenersSendEvent[eName] {
		c.Command(command, event)
	}
	return
}

func (c *commandBus) RegisterCommand(command Command) CommandSettings {
	functionName := getCommandName(command)
	c.listeners[functionName] = []eventName{}
	c.canSend[functionName] = map[eventName]bool{}
	cs := commandSettings{bus: c, command: command, commandName: functionName}

	return &cs
}

func (c *commandBus) Visualize() string {
	s := func(str commandName) string {
		r := strings.NewReplacer("(", "", ")", "", "*", "")
		return r.Replace(strings.TrimPrefix(string(str), "chaturbate/src/services."))
	}
	result := []string{"@startuml"}
	for commandName, _ := range c.listeners {
		if events, ok := c.canSend[commandName]; ok {
			for eventName := range events {
				result = append(result, fmt.Sprintf("%s: %s", s(commandName), eventName))
			}
		}
	}
	for commandName, events := range c.canSend {
		for eventName, _ := range events {
			if commandsNext, ok := c.listenersSendEvent[eventName]; ok {
				for _, commandNext := range commandsNext {
					result = append(result, fmt.Sprintf("%s --> %s: %s", s(commandName), s(getCommandName(commandNext)), eventName))
				}
			} else {
				result = append(result, fmt.Sprintf("%s --> [*]: %s", s(commandName), eventName))
			}
		}
	}
	result = append(result, "@enduml")
	return strings.Join(result, "\n")
}
