package channel

import (
	"errors"
	"fmt"
	"io"

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

func (c *Channel) Broadcast(publisher *client.Client) error {
	if publisher == nil {
		return errors.New("publisher is required")
	}

	for {
		buf := make([]byte, 4096)
		n, err := publisher.Connection.Read(buf)
		if err != nil {
			return err
		}

		for _, cl := range c.Clients {
			for {
				conn := cl.Connection
				_, err := conn.Write(buf[:n])

				if err != nil {
					if err == io.EOF {
						fmt.Printf("receive file complete. \n")
						break
					} else {
						fmt.Printf("conn.read() method execution error, error is:% v \n", err)
						break
					}
				}
			}

		}
	}

}
