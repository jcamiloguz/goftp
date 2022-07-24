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
	FILE
	SUB
	UNSUB
	OK
	ERR
)

type Action struct {
	Id            ACTIONID
	Client        *Client
	Args          map[string]string
	Payload       []byte
	PayloadLength int
}

func NewAction(message []byte, client *Client) (*Action, error) {
	cmd := bytes.ToLower(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))

	// get args from message and convert to map
	args := make(map[string]string)
	actionId, err := GetActionId(string(cmd))
	for _, arg := range bytes.Split(message, []byte(" "))[1:] {
		if bytes.Contains(arg, []byte("=")) {
			key := bytes.Split(arg, []byte("="))[0]
			value := bytes.Split(arg, []byte("="))[1]
			value = bytes.TrimSpace(value)
			args[string(key)] = string(value)
		}
	}
	var payload []byte
	if actionId == FILE {
		content := bytes.Split(message, []byte(" "))[1:]
		payload = bytes.Join(content, []byte(" "))
	} else {
		payload = nil
	}

	if err != nil {
		return nil, err
	}

	return &Action{
		Id:            actionId,
		Client:        client,
		Args:          args,
		Payload:       payload,
		PayloadLength: len(payload),
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
	case "file":
		return FILE, nil
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
	case FILE:
		return "file"
	case OK:
		return "ok"
	case ERR:
		return "error"
	default:
		return "unknown action"
	}
}
