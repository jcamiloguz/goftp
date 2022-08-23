package server

import (
	"errors"
	"fmt"
	"strconv"

	ch "github.com/jcamiloguz/goftp/internal/channel"
	cl "github.com/jcamiloguz/goftp/internal/client"
	"github.com/jcamiloguz/goftp/internal/model"
)

type Config struct {
	Host      string
	Port      string
	NChannels int
}

type Server struct {
	Config *Config

	Channels map[int]*ch.Channel
	Clients  map[string]*cl.Client

	Requests chan *cl.Action
	Response chan *cl.Action

	Outbound chan *model.Payload
	Inbound  chan any
	Payload  *model.Payload
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
		Outbound: make(chan *model.Payload),
		Inbound:  make(chan any),

		Payload: nil,
	}, nil
}

func (s *Server) Start() {
	s.Payload = s.createPayload()
	fmt.Printf("GOFTP Server started on %s:%s\n", s.Config.Host, s.Config.Port)
	fmt.Printf("Number of channels: %d\n", len(s.Channels))
	for {
		select {
		case actionToExc := <-s.Requests:
			isLogged := s.isLogged(actionToExc.Client)
			if !isLogged && actionToExc.Id != cl.REG {
				actionToExc.Client.SendError(errors.New("you are not logged"))
				continue
			}

			actionText := cl.GetActionText(actionToExc.Id)
			fmt.Printf("%s:%s\n", actionToExc.Client.Id, actionText)

			switch actionToExc.Id {

			case cl.REG:
				err := s.register(actionToExc.Client)
				if err != nil {
					actionToExc.Client.SendError(err)
					continue
				}
				actionToExc.Client.SendSuccesful()

			case cl.OUT:
				s.logout(actionToExc.Client)

			case cl.PUB:
				err := s.publish(actionToExc.Client, actionToExc.Args)
				if err != nil {
					actionToExc.Client.SendError(err)
					continue
				}
				actionToExc.Client.SendSuccesful()

			case cl.SUB:
				err := s.subscribe(actionToExc.Client, actionToExc.Args)
				if err != nil {
					actionToExc.Client.SendError(err)
					continue
				}
				actionToExc.Client.SendSuccesful()

			case cl.UNSUB:
				s.unsubscribe(actionToExc.Client)
				actionToExc.Client.SendSuccesful()

			case cl.ERR:
				fmt.Println("Error: ", actionToExc.Args["msg"])
			}
		case <-s.Inbound:
			s.Outbound <- s.Payload
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
		if _, exists := channel.Clients[clientToLogout.Id]; exists {
			channel.RemoveClient(clientToLogout)
			s.Payload.RemoveSubscriber(channel.Id, clientToLogout.Id)
			s.Outbound <- s.Payload
		}
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
		s.Channels[channelId].AddClient(clientToSubscribe)
		s.Payload.AddSubscriber(s.Channels[channelId].Id, clientToSubscribe.Id)
		s.Outbound <- s.Payload
	} else {
		return errors.New("channel does not exist")
	}
	return nil
}

func (s *Server) unsubscribe(clientToUnsubscribe *cl.Client) {
	for _, channel := range s.Channels {
		if _, exists := channel.Clients[clientToUnsubscribe.Id]; exists {
			channel.RemoveClient(clientToUnsubscribe)
			s.Payload.RemoveSubscriber(channel.Id, clientToUnsubscribe.Id)
			s.Outbound <- s.Payload
		}

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
		err := s.Channels[channelId].Broadcast(publisher, file)
		if err != nil {
			return errors.New("error publishing file:" + err.Error())
		}
		s.Payload.AddFile(s.Channels[channelId].Id, model.File{
			Name: fileName,
			Size: size,
		})
		s.Outbound <- s.Payload

	} else {
		return errors.New("error channel does not exist")
	}

	return nil
}
