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
	Response chan *client.Action
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
		Response: make(chan *client.Action),
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
		actionText := client.GetActionText(actionToExc.Id)
		fmt.Println("Action: ", actionText)
		switch actionToExc.Id {
		case client.REG:
			s.register(actionToExc.Client)
		case client.OUT:
			s.logout(actionToExc.Client)
		case client.PUB:
			fmt.Println("Args: ", actionToExc.Args)
			channelToPublish := actionToExc.Args["channel"]
			if channelToPublish == "" {
				actionToExc.Client.Connection.Write([]byte("ERR error channel is required\n"))
				break
			}
			fileName := actionToExc.Args["fileName"]
			if fileName == "" {
				actionToExc.Client.Connection.Write([]byte("ERR error file is required\n"))
				break
			}
			sizeRaw := actionToExc.Args["size"]
			if sizeRaw == "" {
				actionToExc.Client.Connection.Write([]byte("ERR error file is required\n"))
				break
			}
			size, err := strconv.Atoi(sizeRaw)
			if err != nil {
				actionToExc.Client.Connection.Write([]byte("ERR error size is required\n"))
				break
			}
			file := channel.NewFile(fileName, size)
			s.publish(actionToExc.Client, channelToPublish, file)

		case client.SUB:
			// log args
			fmt.Println("Args: ", actionToExc.Args)
			channel := actionToExc.Args["channel"]
			if channel == "" {
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
		s.SendError(newClient, errors.New("you are already logged"))
	} else {
		s.Clients[newClient.Id] = newClient
		s.SendSuccesful(newClient)
	}
}

func (s *Server) logout(clientToLogout *client.Client) {
	if _, exists := s.Clients[clientToLogout.Id]; exists {
		delete(s.Clients, clientToLogout.Id)
		for _, channel := range s.Channels {
			delete(channel.Clients, clientToLogout.Id)
		}
	}
}

func (s *Server) subscribe(clientToSubscribe *client.Client, channel string) {
	channelId, err := strconv.Atoi(channel)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	if _, exists := s.Clients[clientToSubscribe.Id]; exists {
		if _, exists := s.Channels[channelId]; exists {
			s.Channels[channelId].Clients[clientToSubscribe.Id] = clientToSubscribe
			s.SendSuccesful(clientToSubscribe)
		} else {
			s.SendError(clientToSubscribe, errors.New("error channel does not exist"))
		}
	} else {
		s.SendError(clientToSubscribe, errors.New("error you are not logged"))

	}
}

func (s *Server) unsubscribe(clientToUnsubscribe *client.Client) {
	if _, exists := s.Clients[clientToUnsubscribe.Id]; exists {
		for _, channel := range s.Channels {
			delete(channel.Clients, clientToUnsubscribe.Id)
		}
	}
}

func (s *Server) publish(publisher *client.Client, channel string, file *channel.File) {
	channelId, err := strconv.Atoi(channel)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	if _, exists := s.Clients[publisher.Id]; exists {
		if _, exists := s.Channels[channelId]; exists {
			s.SendSuccesful(publisher)
			err := s.Channels[channelId].Broadcast(publisher, file)
			if err != nil {
				s.SendError(publisher, err)
			}
			s.SendSuccesful(publisher)
			s.CleanChannel(channelId)
		} else {
			s.SendError(publisher, errors.New("error channel does not exist"))
		}
	} else {
		s.SendError(publisher, errors.New("error you are not logged"))
	}
}

func (s *Server) SendSuccesful(c *client.Client) {
	okCmd := []byte("OK \n")
	c.Connection.Write(okCmd)
	okAction, err := client.NewAction(okCmd, c)
	if err != nil {
		fmt.Println(err)
	}
	c.Response <- okAction
}

func (s *Server) CleanChannel(channelId int) {
	for _, client := range s.Channels[channelId].Clients {
		err := client.Connection.Close()
		if err != nil {
			fmt.Println(err)
		}
		// delete client
		delete(s.Clients, client.Id)
		delete(s.Channels[channelId].Clients, client.Id)
	}
}

func (s *Server) SendError(c *client.Client, err error) {

	errorMsg := fmt.Sprintf("ERR msg=%s\n", err.Error())
	errorCmd := []byte(errorMsg)
	c.Connection.Write(errorCmd)
	errorAction, err := client.NewAction(errorCmd, c)
	if err != nil {
		fmt.Println(err)
	}
	c.Response <- errorAction
}
