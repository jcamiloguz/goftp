package server

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	ch "github.com/jcamiloguz/goftp/internal/channel"
	cl "github.com/jcamiloguz/goftp/internal/client"
)

type Config struct {
	Host      string
	Port      string
	NChannels int
}

type Server struct {
	Config   *Config
	Channels map[int]*ch.Channel
	Clients  map[string]*cl.Client
	Requests chan *cl.Action
	Response chan *cl.Action
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
		Clients:  make(map[string]*cl.Client),
		Requests: make(chan *cl.Action),
		Response: make(chan *cl.Action),
	}, nil
}

func CreateChannels(NChannels int) map[int]*ch.Channel {
	channels := make(map[int]*ch.Channel)
	for i := 1; i <= NChannels; i++ {
		channel, err := ch.NewChannel(i)
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
		actionToExc := <-s.Requests
		isLogged := s.isLogged(actionToExc.Client)
		if !isLogged && actionToExc.Id != cl.REG {
			s.SendError(actionToExc.Client, errors.New("you are not logged"))
		}
		actionText := cl.GetActionText(actionToExc.Id)
		fmt.Println("Action: ", actionText)

		switch actionToExc.Id {

		case cl.REG:
			err := s.register(actionToExc.Client)
			if err != nil {
				s.SendError(actionToExc.Client, err)
				continue
			}
			s.SendSuccesful(actionToExc.Client)

		case cl.OUT:
			s.logout(actionToExc.Client)

		case cl.PUB:
			err := s.publish(actionToExc.Client, actionToExc.Args)
			if err != nil {
				s.SendError(actionToExc.Client, err)
				continue
			}
			s.SendSuccesful(actionToExc.Client)

		case cl.SUB:
			err := s.subscribe(actionToExc.Client, actionToExc.Args)
			if err != nil {
				s.SendError(actionToExc.Client, err)
				continue
			}
			s.SendSuccesful(actionToExc.Client)

		case cl.UNSUB:
			s.unsubscribe(actionToExc.Client)
			s.SendSuccesful(actionToExc.Client)

		case cl.ERR:
			fmt.Println("Error: ", actionToExc.Args["msg"])
		}
	}

}
func (s *Server) register(newClient *cl.Client) error {
	if _, exists := s.Clients[newClient.Id]; exists {
		return errors.New("error client already registered")
	} else {
		s.Clients[newClient.Id] = newClient
		return nil
	}
}

func (s *Server) logout(clientToLogout *cl.Client) {
	delete(s.Clients, clientToLogout.Id)
	for _, channel := range s.Channels {
		delete(channel.Clients, clientToLogout.Id)
	}
}

func (s *Server) subscribe(clientToSubscribe *cl.Client, args map[string]string) error {
	fmt.Println("Args: ", args)
	channel := args["channel"]
	if channel == "" {
		return errors.New("channel is required")
	}
	channelId, err := strconv.Atoi(channel)
	if err != nil {
		return errors.New("channel must be a number")
	}

	if _, exists := s.Channels[channelId]; exists {
		s.Channels[channelId].Clients[clientToSubscribe.Id] = clientToSubscribe
	} else {
		return errors.New("channel does not exist")
	}
	return nil
}

func (s *Server) unsubscribe(clientToUnsubscribe *cl.Client) {
	for _, channel := range s.Channels {
		delete(channel.Clients, clientToUnsubscribe.Id)
	}

}

func (s *Server) publish(publisher *cl.Client, args map[string]string) error {
	channelToPublish := args["channel"]
	if channelToPublish == "" {
		return errors.New("error channel is required")
	}
	fileName := args["fileName"]
	if fileName == "" {
		return errors.New("error fileName is required")
	}
	sizeRaw := args["size"]
	if sizeRaw == "" {
		return errors.New("error size is required")
	}
	size, err := strconv.Atoi(sizeRaw)
	if err != nil {
		return errors.New("error size is not a number")
	}
	file, err := ch.NewFile(fileName, size)
	if err != nil {
		return err
	}
	channelId, err := strconv.Atoi(channelToPublish)
	if err != nil {
		return errors.New("error channel is not a number")
	}

	if _, exists := s.Channels[channelId]; exists {
		s.SendSuccesful(publisher)
		err := s.Channels[channelId].Broadcast(publisher, file)
		if err != nil {
			return err
		}
		s.CleanChannel(channelId)
	} else {
		return errors.New("error channel does not exist")
	}

	return nil
}

func (s *Server) SendSuccesful(c *cl.Client) {
	okCmd := []byte("OK \n")
	c.Connection.Write(okCmd)
	okAction, err := cl.NewAction(okCmd, c)
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
		delete(s.Clients, client.Id)
		delete(s.Channels[channelId].Clients, client.Id)
	}
}

func (s *Server) SendError(c *cl.Client, err error) {
	errorMsg := fmt.Sprintf("ERR msg=%s\n", err.Error())
	errorCmd := []byte(errorMsg)
	c.Connection.Write(errorCmd)
	errorAction, err := cl.NewAction(errorCmd, c)
	if err != nil {
		fmt.Println(err)
	}
	c.Response <- errorAction
}
func (s *Server) isLogged(c *cl.Client) bool {
	if _, exists := s.Clients[c.Id]; exists {
		return true
	}
	return false
}
