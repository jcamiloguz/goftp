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
		c.Handle(msg)
	}

}

func (c *Client) Handle(message []byte) {
	action, err := NewAction(message, c)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	c.Request <- action
	if action.Id == PUB {
		//wait fot 2 responses: header and file
		<-c.Response
		<-c.Response
	}

	<-c.Response

}
