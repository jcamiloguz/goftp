package server

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/jcamiloguz/goftp/internal/channel"
	"github.com/jcamiloguz/goftp/internal/client"
)

type Config struct {
	Host      string
	Port      string
	NChannels int
}

type Server struct {
	Config   *Config
	Channels map[int]*channel.Channel
	Clients  map[string]*client.Client
	Actions  chan *client.Action
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
		Clients:  make(map[string]*client.Client),
		Actions:  make(chan *client.Action),
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

func (s *Server) Start() {
	fmt.Println("Server started")
	for {
		actionToExc := <-s.Actions
		fmt.Println("Action: ", actionToExc.Id)
		switch actionToExc.Id {
		case client.REG:
			s.register(actionToExc.Client)

		case client.OUT:
			s.logout(actionToExc.Client)

		case client.PUB:
			channel, err := strconv.Atoi(actionToExc.Args["channel"])
			if err != nil {
				actionToExc.Client.Connection.Write([]byte("ERR error channel is required\n"))
				break
			}
			s.publish(actionToExc.Client, channel, actionToExc.Payload)

		case client.SUB:
			channel, err := strconv.Atoi(actionToExc.Args["channel"])
			if err != nil {
				actionToExc.Client.Connection.Write([]byte("ERR error channel is required\n"))
				break
			}
			s.subscribe(actionToExc.Client, channel)

		case client.UNSUB:
			s.unsubscribe(actionToExc.Client)

		case client.ERR:
			actionToExc.Client.Connection.Write([]byte("ERR error\n"))

		}
	}

}
func (s *Server) register(newClient *client.Client) {
	if _, exists := s.Clients[newClient.Id]; exists {
		newClient.Connection.Write([]byte("ERR error you already logged\n"))
	} else {
		s.Clients[newClient.Id] = newClient
		newClient.Connection.Write([]byte("OK\n"))
	}
}

func (s *Server) logout(clientToLogout *client.Client) {
	if _, exists := s.Clients[clientToLogout.Id]; exists {
		delete(s.Clients, clientToLogout.Username)
		for _, channel := range s.Channels {
			delete(channel.Clients, clientToLogout.Id)
		}
	}
}

func (s *Server) subscribe(clientToSubscribe *client.Client, channelId int) {
	if _, exists := s.Clients[clientToSubscribe.Id]; exists {
		if _, exists := s.Channels[channelId]; exists {
			s.Channels[channelId].Clients[clientToSubscribe.Id] = clientToSubscribe
			clientToSubscribe.Connection.Write([]byte("OK\n"))
		} else {
			clientToSubscribe.Connection.Write([]byte("ERR error channel does not exist\n"))
		}
	} else {
		clientToSubscribe.Connection.Write([]byte("ERR error you are not logged\n"))
	}
}

func (s *Server) unsubscribe(clientToUnsubscribe *client.Client) {
	if _, exists := s.Clients[clientToUnsubscribe.Id]; exists {
		for _, channel := range s.Channels {
			delete(channel.Clients, clientToUnsubscribe.Id)
		}
	}
}

func (s *Server) publish(clientToPublish *client.Client, channelId int, payload []byte) {
	if _, exists := s.Clients[clientToPublish.Id]; exists {
		if _, exists := s.Channels[channelId]; exists {
			s.Channels[channelId].Broadcast(clientToPublish.Username, payload)
		} else {
			clientToPublish.Connection.Write([]byte("ERR error channel does not exist\n"))
		}
	} else {
		clientToPublish.Connection.Write([]byte("ERR error you are not logged\n"))
	}
}
