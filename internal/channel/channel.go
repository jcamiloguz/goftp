package channel

import (
	"errors"

	"github.com/jcamiloguz/goftp/internal/client"
)

type Channel struct {
	Id      int16
	Clients map[int]client.Client
}

func NewChannel(idChannel int) (*Channel, error) {
	if idChannel < 1 {
		return nil, errors.New("idChannel is required")
	}

	return &Channel{
		Id:      int16(idChannel),
		Clients: make(map[int]client.Client),
	}, nil
}
func (c *Channel) broadcast(s string, m []byte) {
	msg := append([]byte(s), ": "...)
	msg = append(msg, m...)
	msg = append(msg, '\n')

	for _, cl := range c.Clients {
		conn := cl.Connection
		conn.Write(msg)
	}
}
