package server

import (
	"errors"

	"github.com/jcamiloguz/goftp/internal/client"
)

type Config struct {
	Host      string
	Port      string
	NChannels int
}

type Server struct {
	Config   *Config
	Channels []Channel
}
type Channel struct {
	Id      int16
	Clients map[int][]chan client.Client
}

func NewServer(config *Config) (*Server, error) {
	if config.Host == "" {
		return nil, errors.New("host  is required")
	}
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.NChannels == 0 {
		return nil, errors.New("the number of channels is required")
	}
	channels := CreateChannels(config.NChannels)
	return &Server{
		Config:   config,
		Channels: channels,
	}, nil
}

func NewChannel(idChannel int) *Channel {
	return &Channel{
		Id:      int16(idChannel),
		Clients: make(map[int][]chan client.Client),
	}
}

func CreateChannels(NChannels int) []Channel {
	var channels []Channel
	for i := 0; i < NChannels; i++ {

		channel := NewChannel(i)

		channels = append(channels, *channel)
	}
	return channels

}
