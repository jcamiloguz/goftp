package channel

import (
	"errors"
	"fmt"

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

func (c *Channel) Broadcast(username string, content []byte) error {

	if username == "" {
		return errors.New("username is required")
	}
	if content == nil {
		return errors.New("content is required")
	}

	msg := []byte(fmt.Sprintf("%s: %s\n", username, content))

	for _, cl := range c.Clients {
		conn := cl.Connection
		conn.Write(msg)
	}
	return nil
}
