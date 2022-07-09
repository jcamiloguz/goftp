package client

import (
	"strings"
)

type ACTIONID int

const (
	REG ACTIONID = iota
	OUT
	PUB
	SUB
	UNSUB
	ERR
)

type Action struct {
	Id      ACTIONID
	Client  *Client
	Args    map[string]string
	Payload []byte
}

func NewAction(actionName string, client *Client, args map[string]string, payload []byte) *Action {
	actionId := GetActionId(actionName)

	return &Action{
		Id:      actionId,
		Client:  client,
		Args:    args,
		Payload: payload,
	}
}

func GetActionId(action string) ACTIONID {
	action = strings.ToLower(action)
	switch action {
	case "register":
		return REG
	case "out":
		return OUT
	case "publish":
		return PUB
	case "subscribe":
		return SUB
	case "unsubscribe":
		return UNSUB
	default:
		return ERR
	}
}
