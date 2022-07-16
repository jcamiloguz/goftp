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
type File struct {
	Name string
	Size int
}

func NewFile(name string, size int) *File {
	return &File{
		Name: name,
		Size: size,
	}
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

func (c *Channel) Broadcast(publisher *client.Client, file *File) error {
	if publisher == nil {
		return errors.New("publisher is required")
	}

	fileHeader := fmt.Sprintf("INFO  fileName=%s size=%d", file.Name, file.Size)
	for _, client := range c.Clients {
		client.Connection.Write([]byte(fileHeader))
		fmt.Printf("sending fileheader to %s\n", client.Id)
	}

	for {
		buf := make([]byte, 1024)
		n, err := publisher.Connection.Read(buf)
		if err != nil {
			return err
		}
		fmt.Printf("read %d bytes from publisher\n", n)
		action, err := client.NewAction(buf[:n], publisher)
		if err == nil {
			fmt.Printf("finish")
			switch action.Id {
			case client.OK:
				c.broadcastSuccessful()
				return nil
			case client.ERR:
				c.broadcastError(errors.New("error received file"))
				return errors.New("error received file")
			}
		} else {
			for _, cl := range c.Clients {
				conn := cl.Connection
				_, err := conn.Write(buf[:n])
				if err != nil {
					fmt.Printf("conn.write() method execution error, error is:% v \n", err)
					break
				}
			}
		}
	}
}

/// broadcastSuccessful sends a message to all clients that the file was sent successfully
func (c *Channel) broadcastSuccessful() {
	for _, client := range c.Clients {
		client.Connection.Write([]byte("OK \n"))
	}
}

func (c *Channel) broadcastError(err error) {
	for _, client := range c.Clients {
		client.Connection.Write([]byte(fmt.Sprintf("ERR msg=%s\n", err.Error())))
	}
}
