package common

import "github.com/mikefaraponov/chatum"

func DoNothing() {}

func NewMessage(message string) *chatum.ClientSideEvent {
	return &chatum.ClientSideEvent{
		Message: message,
	}
}

func NewPongMessage() *chatum.ClientSideEvent {
	return &chatum.ClientSideEvent{
		Type: chatum.EventType_PONG,
	}
}
