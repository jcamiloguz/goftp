package client

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

type Client struct {
	Id            string
	Username      string
	Connection    net.Conn
	StartConnTime time.Time
	EndConnTime   time.Time
	Action        chan *Action
}

func NewClient(conn net.Conn, username string, actions chan *Action) (*Client, error) {
	if conn == nil {
		return nil, errors.New("connection is required")
	}
	if username == "" {
		return nil, errors.New("username is required")
	}
	if actions == nil {
		return nil, errors.New("register is required")
	}

	return &Client{
		Id:            conn.RemoteAddr().String(),
		Username:      username,
		Connection:    conn,
		StartConnTime: time.Now(),
		Action:        actions,
	}, nil
}

func (c *Client) Read() error {

	for {
		msg, err := bufio.NewReader(c.Connection).ReadBytes('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// fmt.Println(string(msg))
		c.Handle(msg)
	}

}

func (c *Client) Handle(message []byte) {
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

	action := NewAction(string(cmd), c, args, payload)

	c.Action <- action

}
