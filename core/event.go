package core

import (
	"fmt"

	"github.com/iwataka/mybot/data"
)

const (
	TwitterEventType = "Twitter"
	SlackEventType   = "Slack"
)

type ReceivedEvent struct {
	Type string
	Name string
	Data interface{}
}

func NewReceivedEvent(typ, name string, data interface{}) ReceivedEvent {
	return ReceivedEvent{typ, name, data}
}

func (e ReceivedEvent) String() string {
	return fmt.Sprintf("%s %s event: %#v", e.Type, e.Name, e.Data)
}

type ActionEvent struct {
	Action data.Action
	Data   interface{}
}

func NewActionEvent(action data.Action, data interface{}) ActionEvent {
	return ActionEvent{action, data}
}

func (e ActionEvent) String() string {
	return fmt.Sprintf("%#v -> %#v", e.Action, e.Data)
}
