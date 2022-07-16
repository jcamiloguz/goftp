package client

import (
	"bytes"
	"errors"
	"strings"
)

type ACTIONID int

const (
	REG ACTIONID = iota
	OUT
	PUB
	SUB
	UNSUB
	INFO
	OK
	ERR
)

type Action struct {
	Id      ACTIONID
	Client  *Client
	Args    map[string]string
	Payload []byte
}

func NewAction(message []byte, client *Client) (*Action, error) {
	cmd := bytes.ToLower(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))

	// get args from message and convert to map
	args := make(map[string]string)
	for _, arg := range bytes.Split(message, []byte(" "))[1:] {
		if bytes.Contains(arg, []byte("=")) {
			key := bytes.Split(arg, []byte("="))[0]
			value := bytes.Split(arg, []byte("="))[1]
			value = bytes.TrimSpace(value)
			args[string(key)] = string(value)
		}
	}
	payload := bytes.TrimSpace(bytes.Split(message, []byte(" "))[len(bytes.Split(message, []byte(" ")))-1])
	actionId, err := GetActionId(string(cmd))
	if err != nil {
		return nil, err
	}

	return &Action{
		Id:      actionId,
		Client:  client,
		Args:    args,
		Payload: payload,
	}, nil
}

func GetActionId(action string) (ACTIONID, error) {
	action = strings.ToLower(action)
	switch action {
	case "register":
		return REG, nil
	case "out":
		return OUT, nil
	case "publish":
		return PUB, nil
	case "subscribe":
		return SUB, nil
	case "unsubscribe":
		return UNSUB, nil
	case "ok":
		return OK, nil
	case "error":
		return ERR, nil
	default:
		return ERR, errors.New("unknown action")
	}
}
func GetActionText(action ACTIONID) string {
	switch action {
	case REG:
		return "register"
	case OUT:
		return "out"
	case PUB:
		return "publish"
	case SUB:
		return "subscribe"
	case UNSUB:
		return "unsubscribe"
	case INFO:
		return "info"
	case OK:
		return "ok"
	case ERR:
		return "error"
	default:
		return "unknown action"
	}
}
