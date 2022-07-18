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
type File struct {
	Name string
	Size int
}

func NewFile(name string, size int) (*File, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if size < 1 {
		return nil, errors.New("size is required")
	}
	return &File{
		Name: name,
		Size: size,
	}, nil
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
	writers := make([]io.Writer, 0, len(c.Clients))
	for _, client := range c.Clients {
		writers = append(writers, client.Connection)
	}
	writer := io.MultiWriter(writers...)

	fileHeader := fmt.Sprintf("INFO  fileName=%s size=%d", file.Name, file.Size)

	_, err := writer.Write([]byte(fileHeader))
	if err != nil {
		return err
	}
	// create a MultiWriter to send the file to all clients
	io.Copy(writer, publisher.Connection)
	return nil
}
