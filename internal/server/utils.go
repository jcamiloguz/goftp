package server

import (
	"fmt"

	cl "github.com/jcamiloguz/goftp/internal/client"
)

func (s *Server) SendSuccesful(c *cl.Client) {
	okCmd := []byte("ok \n")
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
	errorMsg := fmt.Sprintf("error msg=%s\n", err.Error())
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
