package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/jcamiloguz/goftp/internal/channel"
	"github.com/jcamiloguz/goftp/internal/client"
)

type Config struct {
	Host      string
	Port      string
	NChannels int
}

type Server struct {
	Config          *Config
	Channels        map[int]*channel.Channel
	Clients         map[string]*client.Client
	Actions         chan *client.Client
	Login           chan *client.Client
	Logout          chan *client.Client
	Registrations   chan *client.Client
	DeRegistrations chan *client.Client
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

func CreateChannels(NChannels int) map[int]*channel.Channel {
	channels := make(map[int]*channel.Channel)
	for i := 1; i <= NChannels; i++ {

		channel, err := channel.NewChannel(i)
		if err != nil {
			log.Println(err)
		}
		channels[i] = channel
	}
	return channels

}
func (s *Server) handleConnection(conn net.Conn) {

}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handleConnection(conn)
	}

}
