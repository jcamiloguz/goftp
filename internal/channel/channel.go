package channel

import (
	"errors"
	"fmt"
	"io"

	cl "github.com/jcamiloguz/goftp/internal/client"
)

const BUFFER_SIZE = 1024

type Channel struct {
	Id      int16
	Clients map[string]*cl.Client
	Writer  io.Writer
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
		Clients: make(map[string]*cl.Client),
	}, nil
}

func (ch *Channel) AddClient(client *cl.Client) {
	ch.Clients[client.Id] = client
	writers := make([]io.Writer, 0, len(ch.Clients))
	for _, client := range ch.Clients {
		writers = append(writers, client.Connection)
	}

	multiWriter := io.MultiWriter(writers...)

	ch.Writer = multiWriter
}

func (ch *Channel) RemoveClient(client *cl.Client) {
	delete(ch.Clients, client.Id)
	writers := make([]io.Writer, 0, len(ch.Clients))
	for _, client := range ch.Clients {
		writers = append(writers, client.Connection)
	}
	multiWriter := io.MultiWriter(writers...)

	ch.Writer = multiWriter
}

func (ch *Channel) Broadcast(publisher *cl.Client, file *File) error {
	if len(ch.Clients) == 0 {
		return errors.New("no clients in channel")
	}
	if publisher == nil {
		return errors.New("publisher is required")
	}
	if file == nil {
		return errors.New("file is required")
	}
	if len(ch.Clients) < 1 {
		return errors.New("no clients connected")
	}

	fileHeader := fmt.Sprintf("PUB  fileName=%s size=%d ", file.Name, file.Size)
	buffHead := make([]byte, BUFFER_SIZE)
	copy(buffHead, []byte(fileHeader))

	_, err := ch.Writer.Write(buffHead)
	if err != nil {
		return err
	}

	for {
		buff := make([]byte, BUFFER_SIZE)

		_, err := publisher.Connection.Read(buff)
		if err != nil {
			errMsg := fmt.Sprintf("error reading connection from publisher -- %s", err.Error())
			return errors.New(errMsg)
		}

		action, err := cl.NewAction(buff, publisher)
		if err != nil {
			errMsg := fmt.Sprintf("error geting action from publisher except a File/ok action -- %s", err.Error())
			return errors.New(errMsg)

		}

		if action.Id == cl.FILE {
			_, err := ch.Writer.Write(buff)
			if err != nil {
				errMsg := fmt.Sprintf("error writing file -- %s", err.Error())
				return errors.New(errMsg)
			}
		}

		if action.Id == cl.OK {
			_, err := ch.Writer.Write(buff)
			if err != nil {
				errMsg := fmt.Sprintf("error sending ok final confirmation  -- %s", err.Error())
				return errors.New(errMsg)
			}
			break
		}
	}

	return nil
}
