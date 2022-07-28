package server

import (
	"fmt"

	cl "github.com/jcamiloguz/goftp/internal/client"
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
