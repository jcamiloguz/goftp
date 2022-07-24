package channel

import (
	"errors"
	"fmt"
	"io"

	cl "github.com/jcamiloguz/goftp/internal/client"
)

type Channel struct {
	Id      int16
	Clients map[string]*cl.Client
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

func (c *Channel) Broadcast(publisher *cl.Client, file *File) error {
	if publisher == nil {
		return errors.New("publisher is required")
	}

	writers := make([]io.Writer, 0, len(c.Clients))
	for _, client := range c.Clients {
		writers = append(writers, client.Connection)
	}
	writer := io.MultiWriter(writers...)
	fileHeader := fmt.Sprintf("publish  fileName=%s size=%d", file.Name, file.Size)

	_, err := writer.Write([]byte(fileHeader))
	if err != nil {
		return err
	}
	buff := make([]byte, 1024)
	for {
		_, err := publisher.Connection.Read(buff)
		if err != nil {

			errMsg := fmt.Sprintf("error reading connection from publisher -- %s", err.Error())
			return errors.New(errMsg)
		}
		// get action form publisher
		action, err := cl.NewAction(buff, publisher)
		if err != nil {
			errMsg := fmt.Sprintf("error geting action from publisher except a File/ok action -- %s", err.Error())
			return errors.New(errMsg)

		}
		fmt.Printf("%d\n", action.Id)
		if action.Id == cl.FILE {
			n, err := writer.Write(buff)
			fmt.Printf("%d\n", n)
			if err != nil {
				errMsg := fmt.Sprintf("error writing file -- %s", err.Error())
				return errors.New(errMsg)
			}
			continue
		}

		if action.Id == cl.OK {
			_, err := writer.Write(buff)
			if err != nil {
				errMsg := fmt.Sprintf("error sending ok final confirmation  -- %s", err.Error())
				return errors.New(errMsg)
			}
			break
		}
	}
	return nil
}
