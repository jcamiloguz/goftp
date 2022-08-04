package server

import (
	"fmt"
	"log"

	ch "github.com/jcamiloguz/goftp/internal/channel"
	cl "github.com/jcamiloguz/goftp/internal/client"
	"github.com/jcamiloguz/goftp/internal/model"
)

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
func (s *Server) isLogged(c *cl.Client) bool {
	if _, exists := s.Clients[c.Id]; exists {
		return true
	}
	return false
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

func (s *Server) createPayload() *model.Payload {
	payload := &model.Payload{
		Channels: []model.Channel{},
	}
	for _, channel := range s.Channels {
		channelModel := &model.Channel{
			Id:          channel.Id,
			Subscribers: []model.Client{},
			Files:       []model.File{},
		}
		// add channel to payload
		payload.Channels = append(payload.Channels, *channelModel)

	}
	return payload
}
