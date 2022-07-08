package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
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

type Client struct {
	Username      string
	Connection    net.Conn
	StartConnTime time.Time
	EndConnTime   time.Time
	RegisterChn   chan<- *Client
	DeregisterChn chan<- *Client
}

func NewClient(conn net.Conn, username string, registerChn chan<- *Client, deregisterChn chan<- *Client) (*Client, error) {
	if conn == nil {
		return nil, errors.New("connection is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if registerChn == nil {
		return nil, errors.New("register is required")
	}
	if deregisterChn == nil {
		return nil, errors.New("deregister is required")
	}

	return &Client{
		Username:      username,
		Connection:    conn,
		StartConnTime: time.Now(),
		RegisterChn:   registerChn,
		DeregisterChn: deregisterChn,
	}, nil
}

func (c *Client) Read() error {

	for {
		msg, err := bufio.NewReader(c.Connection).ReadBytes('\n')
		if err == io.EOF {
			// Cierra la connetion
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println(string(msg))
		c.Handle(msg)
	}

}

func (c *Client) Handle(message []byte) {
	cmd := bytes.ToLower(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))
	action := GetActionId(string(cmd))
	switch action {
	case REG:
		c.Register(args)
	case OUT:
		c.Out(args)
	case PUB:
		c.Publish(args)
	case SUB:
		c.Subscribe(args)
	case UNSUB:
		c.Unsubscribe(args)
	default:
		fmt.Println(string(cmd))
	}
}

func (c *Client) Register(args []byte) error {
	fmt.Println("Registering Client")
	return nil
}

func (c *Client) Publish(args []byte) error {
	fmt.Println("Publishing Message")
	return nil
}

func (c *Client) Subscribe(args []byte) error {
	fmt.Println("Publishing Client")
	return nil
}

func (c *Client) Unsubscribe(args []byte) error {
	fmt.Println("Publishing Client")
	return nil
}

func (c *Client) Out(args []byte) error {
	return nil
}

func GetActionId(action string) ACTIONID {

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
