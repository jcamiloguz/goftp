package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	Id            string
	Connection    net.Conn
	StartConnTime time.Time
	EndConnTime   time.Time
	Request       chan *Action
	Response      chan *Action
}

func NewClient(conn net.Conn, actions chan *Action, responses chan *Action) (*Client, error) {
	if conn == nil {
		return nil, errors.New("connection is required")
	}
	if actions == nil {
		return nil, errors.New("register is required")
	}
	if responses == nil {
		return nil, errors.New("response is required")
	}

	return &Client{
		Id:            conn.RemoteAddr().String(),
		Connection:    conn,
		StartConnTime: time.Now(),
		Request:       actions,
		Response:      responses,
	}, nil
}

func (c *Client) Read() error {
	fmt.Printf("client: %v\n", c.Id)

	for {
		msg, err := bufio.NewReader(c.Connection).ReadBytes('\n')
		if err == io.EOF {
			outAction, err := NewAction([]byte("out"), c)
			if err != nil {
				return err
			}
			c.Request <- outAction
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("message: %s\n", msg)
		err = c.Handle(msg)
		if err != nil {
			return err
		}
	}

}

func (c *Client) Handle(message []byte) error {
	action, err := NewAction(message, c)
	if err != nil {
		fmt.Printf("handle err: ")
		return err
	}
	fmt.Printf("action: %v\n", action.Id)

	c.Request <- action

	response := <-c.Response
	if response.Id == ERR {
		return errors.New(response.Args["msg"])
	}
	return nil
}
