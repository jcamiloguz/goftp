package channel

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/jcamiloguz/goftp/internal/client"
)

type Channel struct {
	Id      int16
	Clients map[string]*client.Client
}

func NewChannel(idChannel int) (*Channel, error) {
	if idChannel < 1 {
		return nil, errors.New("idChannel is required")
	}

	return &Channel{
		Id:      int16(idChannel),
		Clients: make(map[string]*client.Client),
	}, nil
}

func (c *Channel) Broadcast(sender *client.Client, content []byte) error {

	if sender == nil {
		return errors.New("sender is required")
	}
	if content == nil {
		return errors.New("content is required")
	}

	msg := []byte(fmt.Sprintf("%s: %s\n", sender.Username, content))

	for _, cl := range c.Clients {
		conn := cl.Connection
		_, err := conn.Write(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Channel) SendFile(publisher net.Conn, Subscriber net.Conn) error {
	if publisher == nil {
		return errors.New("publisher is required")
	}
	if Subscriber == nil {
		return errors.New("Subscriber is required")
	}

	for {
		msg, err := bufio.NewReader(publisher).ReadBytes('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = Subscriber.Write(msg)
		if err != nil {
			return err
		}
	}
}
