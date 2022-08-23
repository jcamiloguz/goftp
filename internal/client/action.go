package client

import (
	"bytes"
	"errors"
	"fmt"
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
			cmdQuery := bytes.Split(arg, []byte("="))
			key := cmdQuery[0]
			value := cmdQuery[1]
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
	action = strings.ToUpper(action)
	switch action {
	case "REG":
		return REG, nil
	case "OUT":
		return OUT, nil
	case "PUB":
		return PUB, nil
	case "FILE":
		return FILE, nil
	case "SUB":
		return SUB, nil
	case "UNSUB":
		return UNSUB, nil
	case "OK":
		return OK, nil
	case "ERR":
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

func (c *Client) SendSuccesful() {
	okCmd := []byte("OK \n")
	c.Connection.Write(okCmd)
	okAction, err := NewAction(okCmd, c)
	if err != nil {
		fmt.Println(err)
	}
	c.Response <- okAction
}

func (c *Client) SendError(err error) {
	errMsg := strings.Replace(err.Error(), " ", "_", -1)
	errorMsg := fmt.Sprintf("ERR msg=\"%s\"\n", errMsg)
	errorCmd := []byte(errorMsg)
	c.Connection.Write(errorCmd)
	errorAction, err := NewAction(errorCmd, c)
	if err != nil {
		fmt.Println(err)
	}
	c.Response <- errorAction
}

func (c *Client) SendPublishFileHeader(fileName string, size int) {
	fileHeader := fmt.Sprintf("PUB  fileName=%s size=%d ", fileName, size)
	fileHeaderCmd := []byte(fileHeader)
	c.Connection.Write(fileHeaderCmd)
}
